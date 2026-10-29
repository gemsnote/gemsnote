"""Celery SMTP worker. Broker task arguments contain IDs, never credentials."""
import json
import os
import smtplib
import ssl
from email.message import EmailMessage
from email.utils import formatdate

from celery import Celery
from celery.utils.log import get_task_logger
from redis import Redis

logger = get_task_logger(__name__)

REDIS_URL = os.environ['GEMSNOTE_MAIL_REDIS_URL']
app = Celery('gemsnote-mail', broker=REDIS_URL)
app.conf.update(
    task_default_queue='gemsnote.mail', accept_content=['json'],
    task_serializer='json', task_ignore_result=True,
    worker_prefetch_multiplier=1, task_acks_late=True,
    task_soft_time_limit=60, task_time_limit=90,
    broker_connection_retry_on_startup=True,
)
store = Redis.from_url(REDIS_URL, socket_connect_timeout=3, socket_timeout=3)


def status(task_id, state, error=None):
    value = {'State': state}
    if error:
        value['Error'] = error
    store.setex('gemsnote:mail:status:' + task_id, 86400, json.dumps(value))


def deliver(payload, task_id):
    settings = payload['smtp']
    message = EmailMessage()
    message['From'] = settings['username']
    message['To'] = ', '.join(payload['to'])
    message['Subject'] = payload['subject']
    message['Date'] = formatdate(localtime=False)
    message['Message-ID'] = f'<{task_id}@gemsnote>'
    message.set_content('Please view this message in an HTML-capable email client.')
    message.add_alternative(payload['body'], subtype='html')
    context = ssl.create_default_context()
    if settings['ssl']:
        client = smtplib.SMTP_SSL(settings['host'], settings['port'], timeout=20, context=context)
    else:
        client = smtplib.SMTP(settings['host'], settings['port'], timeout=20)
    with client:
        client.ehlo()
        if not settings['ssl']:
            # Never transmit SMTP credentials over an unencrypted connection.
            client.starttls(context=context)
            client.ehlo()
        client.login(settings['username'], settings['password'])
        refused = client.send_message(message)
        if refused:
            raise smtplib.SMTPRecipientsRefused(refused)


@app.task(name='gemsnote.mail.send')
def send(task_id):
    key = 'gemsnote:mail:payload:' + task_id
    previous = store.get('gemsnote:mail:status:' + task_id)
    if previous and json.loads(previous).get('State') in ('sent', 'failed'):
        return
    raw = store.get(key)
    if raw is None:
        status(task_id, 'failed', 'mailTaskExpired')
        return
    status(task_id, 'sending')
    error = None
    try:
        deliver(json.loads(raw), task_id)
    except smtplib.SMTPAuthenticationError:
        error = 'mailAuthenticationFailed'
    except smtplib.SMTPRecipientsRefused:
        error = 'mailRecipientRejected'
    except (TimeoutError, OSError):
        error = 'mailConnectionFailed'
    except Exception:
        # Do not expose SMTP server responses, passwords, tokens or message bodies.
        error = 'mailDeliveryFailed'
    if error:
        logger.warning('Mail task %s failed: %s', task_id, error)
    status(task_id, 'failed' if error else 'sent', error)
    store.delete(key)
