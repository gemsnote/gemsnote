-- Gemsnote PostgreSQL installation seed.
--
-- Generated from mongodb_backup/gemsnote_install_data with the Gemsnote
-- MongoDB-to-PostgreSQL migration tool. Runtime-only collections are omitted:
-- sessions, tokens, email_logs, suggestions, and reports.
-- Every INSERT is idempotent so this file can safely be applied more than once.

-- Dumped from database version 15.19 (Debian 15.19-1.pgdg13+2)
-- Dumped by pg_dump version 15.19 (Debian 15.19-1.pgdg13+2)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

BEGIN;

--
-- Data for Name: albums; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.albums (id, user_id, name, type, seq, created_time) VALUES ('53aa5f9e99c37b4525000001', '52d3e8ac99c37b7f0d000001', 'dddeeee', 0, -1, '2014-06-25 05:35:26.007+00') ON CONFLICT DO NOTHING;
INSERT INTO public.albums (id, user_id, name, type, seq, created_time) VALUES ('53aa7c2c99c37b4525000002', '52d3e8ac99c37b7f0d000001', 'xxx', 0, -1, '2014-06-25 07:37:16.148+00') ON CONFLICT DO NOTHING;
INSERT INTO public.albums (id, user_id, name, type, seq, created_time) VALUES ('53aa7c6a99c37b4525000003', '52d3e8ac99c37b7f0d000001', 'Beauty', 0, -1, '2014-06-25 07:38:18.992+00') ON CONFLICT DO NOTHING;
INSERT INTO public.albums (id, user_id, name, type, seq, created_time) VALUES ('53aa7e4399c37b5a8e000001', '52d3e8ac99c37b7f0d000001', 'you ', 0, -1, '2014-06-25 07:46:11.768+00') ON CONFLICT DO NOTHING;
INSERT INTO public.albums (id, user_id, name, type, seq, created_time) VALUES ('53aa7fe399c37b5e6f000001', '52d3e8ac99c37b7f0d000001', 'you can make it!', 0, -1, '2014-06-25 07:53:07.497+00') ON CONFLICT DO NOTHING;
INSERT INTO public.albums (id, user_id, name, type, seq, created_time) VALUES ('53aed9a499c37b290f000001', '52d8d30799c37b1fb2000001', 'hello', 0, -1, '2014-06-28 15:05:08.406+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: attachs; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541e722a99c37b5dc0000002', '541e6e7b19807a705c000000', '52d8d30799c37b1fb2000001', '5b5cb156f1e98f9b939dec65524a74e2.html', 'a.html', 586, 'html', 'files/52d8d30799c37b1fb2000001/attachs', '2014-09-21 06:37:30.216+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541e748999c37b60a9000001', '541e6e7b19807a705c000000', '52d8d30799c37b1fb2000001', 'bdc841a9a73b49d0446b93bb0480d15a.pdf', 'dont match twice.pdf', 936461, 'pdf', 'files/52d8d30799c37b1fb2000001/attachs', '2014-09-21 06:47:37.294+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541e754299c37b60a9000002', '541e6e7b19807a705c000000', '52d8d30799c37b1fb2000001', 'b08067378ebf483059439b0f86d7f02a.html', 'a.html', 586, 'html', 'files/52d8d30799c37b1fb2000001/attachs', '2014-09-21 06:50:42.534+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541e755399c37b60a9000003', '541e6e7b19807a705c000000', '52d8d30799c37b1fb2000001', 'd32cd183d137c08fe159035593b0b6ac.doc', 'DPCS2014-registry.doc', 163328, 'doc', 'files/52d8d30799c37b1fb2000001/attachs', '2014-09-21 06:50:59.466+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541e75ed99c37b644b000001', '541e6e7b19807a705c000000', '52d8d30799c37b1fb2000001', 'ce811efe9cfc8704a2c338c324205d83.css', 'bootstrap.3.2.0.min.css', 109518, 'css', 'files/52d8d30799c37b1fb2000001/attachs', '2014-09-21 06:53:33.703+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541e760799c37b644b000002', '541e6e7b19807a705c000000', '52d8d30799c37b1fb2000001', '73455fd1ccb3a44c94d33ebdb2bd7bf0.pdf', 'dont match twice.pdf', 936461, 'pdf', 'files/52d8d30799c37b1fb2000001/attachs', '2014-09-21 06:53:59.02+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541e763f99c37b64e0000001', '541e6e7b19807a705c000000', '52d8d30799c37b1fb2000001', '47e0a541198c1dcab1575377c1ce16d1.html', 'a.html', 586, 'html', 'files/52d8d30799c37b1fb2000001/attachs', '2014-09-21 06:54:55.391+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541e766699c37b653d000001', '541e6e7b19807a705c000000', '52d8d30799c37b1fb2000001', '0be119973faf414114cc9033d13f9dc3.pdf', 'dont match twice.pdf', 936461, 'pdf', 'files/52d8d30799c37b1fb2000001/attachs', '2014-09-21 06:55:34.63+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541e76a899c37b65af000001', '541e6e7b19807a705c000000', '52d8d30799c37b1fb2000001', '1f3132fffcfe702acb9ee3ec61fdfe24.js', 'excanvas.js', 19314, 'js', 'files/52d8d30799c37b1fb2000001/attachs', '2014-09-21 06:56:40.667+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541e76da99c37b6626000001', '541e6e7b19807a705c000000', '52d8d30799c37b1fb2000001', '93c634170f95b3446ccb659563d30b8e.pdf', 'dont match twice.pdf', 936461, 'pdf', 'files/52d8d30799c37b1fb2000001/attachs', '2014-09-21 06:57:30.69+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541e770599c37b666b000001', '541e6e7b19807a705c000000', '52d8d30799c37b1fb2000001', '2af02c272ddedd76865e6b90f0dbe902.pdf', 'dont match twice.pdf', 936461, 'pdf', 'files/52d8d30799c37b1fb2000001/attachs', '2014-09-21 06:58:13.594+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541e778999c37b66b1000001', '541e6e7b19807a705c000000', '52d8d30799c37b1fb2000001', '5f3b32f681960bee20165d15556354ba.css', 'bootstrap.3.2.0.min.css', 109518, 'css', 'files/52d8d30799c37b1fb2000001/attachs', '2014-09-21 07:00:25.795+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541e781d99c37b67c7000001', '541e6e7b19807a705c000000', '52d8d30799c37b1fb2000001', 'c2adf8f8558929b952f9525ec98ae15a.html', 'a.html', 586, 'html', 'files/52d8d30799c37b1fb2000001/attachs', '2014-09-21 07:02:53.641+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541e782799c37b67c7000002', '541e6e7b19807a705c000000', '52d8d30799c37b1fb2000001', '7473694513a5d13a6a286c9cd1d2eb50.html', 'a.html', 586, 'html', 'files/52d8d30799c37b1fb2000001/attachs', '2014-09-21 07:03:03.86+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541e9d1599c37b798e00000c', '541d2a7f99c37b1947000001', '52d8d30799c37b1fb2000001', 'a4fe5c4732324364915fdc2fc6830fde.pdf', 'dont match twice.pdf', 936461, 'pdf', 'files/52d8d30799c37b1fb2000001/attachs/a4fe5c4732324364915fdc2fc6830fde.pdf', '2014-09-21 09:40:37.486+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541ea56c99c37ba33b000001', '541d2a7f99c37b1947000001', '52d8d30799c37b1fb2000001', 'd0098bc6f940aea82a70c962c2a0785f.html', 'a.html', 586, 'html', 'files/52d8d30799c37b1fb2000001/attachs/d0098bc6f940aea82a70c962c2a0785f.html', '2014-09-21 10:16:12.221+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541ea59199c37ba33b000002', '541d2a7f99c37b1947000001', '52d8d30799c37b1fb2000001', '84a3d581c2630136105b8169308f814c.doc', '101310127_朱斌_毕业论文.doc', 1026560, 'doc', 'files/52d8d30799c37b1fb2000001/attachs/84a3d581c2630136105b8169308f814c.doc', '2014-09-21 10:16:49.091+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541ea5ed99c37ba33b000003', '541d23c399c37b0eb5000001', '52d8d30799c37b1fb2000001', '5a7efc8dc1afa7f63a55fac6f6ff911b.doc', '101310127_朱斌_毕业论文.doc', 1026560, 'doc', 'files/52d8d30799c37b1fb2000001/attachs/5a7efc8dc1afa7f63a55fac6f6ff911b.doc', '2014-09-21 10:18:21.759+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541ea5f399c37ba33b000004', '541d23c399c37b0eb5000001', '52d8d30799c37b1fb2000001', '9419c95f83a58f8b669c0124323864b6.doc', 'DPCS2014-registry.doc', 163328, 'doc', 'files/52d8d30799c37b1fb2000001/attachs/9419c95f83a58f8b669c0124323864b6.doc', '2014-09-21 10:18:27.965+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541ea5f799c37ba33b000005', '541d23c399c37b0eb5000001', '52d8d30799c37b1fb2000001', 'f729a80534de32a029e8c3a5cd6e7b72.js', 'excanvas.js', 19314, 'js', 'files/52d8d30799c37b1fb2000001/attachs/f729a80534de32a029e8c3a5cd6e7b72.js', '2014-09-21 10:18:31.092+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541eae2199c37ba33b000007', '541d261799c37b11be000002', '52d8d30799c37b1fb2000001', '90d43d3ed2bd8cb18a60f1772e84926a.css', 'bootstrap.3.2.0.min.css', 109518, 'css', 'files/52d8d30799c37b1fb2000001/attachs/90d43d3ed2bd8cb18a60f1772e84926a.css', '2014-09-21 10:53:21.876+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541eb2c499c37bacd9000001', '541eb1ca19807a1650000000', '52d8d30799c37b1fb2000001', '08f21209e2af570663eb45910f9bd983.doc', '101310127_朱斌_毕业论文.doc', 1026560, 'doc', 'files/52d8d30799c37b1fb2000001/attachs/08f21209e2af570663eb45910f9bd983.doc', '2014-09-21 11:13:08.77+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541eb2c499c37bacd9000003', '541eb1ca19807a1650000000', '52d8d30799c37b1fb2000001', 'c8bc335eb4ad8b2657f0b2e269f12209.css', 'bootstrap.3.2.0.min.css', 109518, 'css', 'files/52d8d30799c37b1fb2000001/attachs/c8bc335eb4ad8b2657f0b2e269f12209.css', '2014-09-21 11:13:08.771+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541eb52d99c37bb007000002', '541eb1ca19807a1650000000', '52d8d30799c37b1fb2000001', '65adf6e902489600d6e216f46a1261ae.doc', '101310127_朱斌_毕业论文.doc', 1026560, 'doc', 'files/52d8d30799c37b1fb2000001/attachs/65adf6e902489600d6e216f46a1261ae.doc', '2014-09-21 11:23:25.048+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541eb6230716dcbd8f1f6488', '541eb62399c37bb1bd000001', '52d8d30799c37b1fb2000001', 'b8624bdbc57e5d30a4a58283cd133d33.doc', '101310127_朱斌_毕业论文.doc', 1026560, 'doc', 'b8624bdbc57e5d30a4a58283cd133d33.doc', '2014-09-21 11:27:31.826+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541eb6230716dcbd8f1f6489', '541eb62399c37bb1bd000001', '52d8d30799c37b1fb2000001', '2a5403b5da9fcaaf6001fcdda1a44ed5.css', 'bootstrap.3.2.0.min.css', 109518, 'css', '2a5403b5da9fcaaf6001fcdda1a44ed5.css', '2014-09-21 11:27:31.829+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541eb6230716dcbd8f1f648a', '541eb62399c37bb1bd000001', '52d8d30799c37b1fb2000001', 'abb109ce06360db9b78c88e3ae92016d.doc', '101310127_朱斌_毕业论文.doc', 1026560, 'doc', 'abb109ce06360db9b78c88e3ae92016d.doc', '2014-09-21 11:27:31.835+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541eb6810716dcbd8f1f648c', '541eb68199c37bb252000001', '52d8d30799c37b1fb2000001', '25f3a93dca3d3ff0d2897e71d54261e5.doc', '101310127_朱斌_毕业论文.doc', 1026560, 'doc', 'files/52d3e8ac99c37b7f0d000001/attachs/25f3a93dca3d3ff0d2897e71d54261e5.doc', '2014-09-21 11:29:05.599+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541eb6810716dcbd8f1f648d', '541eb68199c37bb252000001', '52d8d30799c37b1fb2000001', 'aa58c310cb73014130b93b987db0865a.css', 'bootstrap.3.2.0.min.css', 109518, 'css', 'files/52d3e8ac99c37b7f0d000001/attachs/aa58c310cb73014130b93b987db0865a.css', '2014-09-21 11:29:05.602+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541eb69499c37bb252000002', '541eb1ca19807a1650000000', '52d8d30799c37b1fb2000001', 'cf26f59d9ca86d277140845be3a5eb1d.css', 'bootstrap.3.2.0.min.css', 109518, 'css', 'files/52d8d30799c37b1fb2000001/attachs/cf26f59d9ca86d277140845be3a5eb1d.css', '2014-09-21 11:29:24.967+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541eb99b99c37bb571000003', '541d25dc09915b19a9000000', '52d3e8ac99c37b7f0d000001', '45a9b414211f69afd98f16b2a3bce04f.jpg', 'oYYBAFN5yriIUJAsAAFAV_Q2ob8AAAeVQMuDMkAAUBv517.jpg', 82007, 'jpg', 'files/52d3e8ac99c37b7f0d000001/attachs/45a9b414211f69afd98f16b2a3bce04f.jpg', '2014-09-21 11:42:19.065+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541ebcfd99c37bb798000002', '541eb10019807a3063000000', '52d8d30799c37b1fb2000001', '224879b6d3ed76999f70781387ace3ee.css', 'bootstrap.3.2.0.min.css', 109518, 'css', 'files/52d8d30799c37b1fb2000001/attachs/224879b6d3ed76999f70781387ace3ee.css', '2014-09-21 11:56:45.334+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541ec1cd99c37bbdd2000002', '541eb10019807a3063000000', '52d8d30799c37b1fb2000001', '9006ae80400309df2f82892999bf21bc.css', 'bootstrap.3.2.0.min.css', 109518, 'css', 'files/52d8d30799c37b1fb2000001/attachs/9006ae80400309df2f82892999bf21bc.css', '2014-09-21 12:17:17.149+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541ed89099c37bc5c6000001', '541ec60e19807a7f65000000', '52d8d30799c37b1fb2000001', '4660fdb3f83acb1c6a548d1d008b117e.doc', '101310127_朱斌_毕业论文.doc', 1026560, 'doc', 'files/52d8d30799c37b1fb2000001/attachs/4660fdb3f83acb1c6a548d1d008b117e.doc', '2014-09-21 13:54:24.754+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541edfa299c37be66b000001', '541eb0e219807a0bfc000000', '52d8d30799c37b1fb2000001', 'e42b495592843d06a1f389dce4efab5e.jpg', 'u=365369868,865426152&fm=23&gp=0.jpg', 23978, 'jpg', 'files/52d8d30799c37b1fb2000001/attachs/e42b495592843d06a1f389dce4efab5e.jpg', '2014-09-21 14:24:34.425+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541ef95899c37b0296000004', '541eb10019807a3063000000', '52d8d30799c37b1fb2000001', '5b1dac38ed16c4e37857b33aa5b0354e.png', 'QQ20140513-1@2x.png', 20083, 'png', 'files/52d8d30799c37b1fb2000001/attachs/5b1dac38ed16c4e37857b33aa5b0354e.png', '2014-09-21 16:14:16.776+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('541efed799c37b0b71000012', '541eb10019807a3063000000', '52d8d30799c37b1fb2000001', '827a63401e318390e3175675eb33bcff.png', 'QQ20140623-1@2x.png', 499836, 'png', 'files/52d8d30799c37b1fb2000001/attachs/827a63401e318390e3175675eb33bcff.png', '2014-09-21 16:37:43.686+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('54200f9199c37b8c81000001', '5420047519807a248b000000', '52d8d30799c37b1fb2000001', 'e20824357faabe4abe09080fde79b1e3.png', 'gemsnote_green.png', 15799, 'png', 'files/52d8d30799c37b1fb2000001/attachs/e20824357faabe4abe09080fde79b1e3.png', '2014-09-22 12:01:21.162+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('54200fb199c37b8c8100000a', '5420047519807a248b000000', '52d8d30799c37b1fb2000001', '9172c14d15ec49b581ad9285a8945a90.psd', 'gemsnote-icon.psd', 46673, 'psd', 'files/52d8d30799c37b1fb2000001/attachs/9172c14d15ec49b581ad9285a8945a90.psd', '2014-09-22 12:01:53.275+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('5420331799c37b90c1000002', '5420326319807a2491000000', '52d8d30799c37b1fb2000001', 'a3cbeddab9e451fa4a3c8a5802e7dcc9.png', 'gemsnote-64.png', 5799, 'png', 'files/52d8d30799c37b1fb2000001/attachs/a3cbeddab9e451fa4a3c8a5802e7dcc9.png', '2014-09-22 14:32:55.175+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('54257b0a99c37bca37000001', '5371aa4e19807a273a000000', '52d26b4e99c37b609a000001', '5b968eeb36a2764aeb7e912777097b6e.html', 'a.html', 586, 'html', 'files/52d26b4e99c37b609a000001/attachs/5b968eeb36a2764aeb7e912777097b6e.html', '2014-09-26 14:41:14.626+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('54257d5199c37bca37000002', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', '5c459345a8831fa94d59389b984879f0.jpg', '20140922_114131.jpg', 2167194, 'jpg', 'files/52d26b4e99c37b609a000001/attachs/5c459345a8831fa94d59389b984879f0.jpg', '2014-09-26 14:50:57.358+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('54257d5199c37bca37000003', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', '768dd69701318651ed288dc0368303e5.doc', '101310127_朱斌_毕业论文.doc', 1026560, 'doc', 'files/52d26b4e99c37b609a000001/attachs/768dd69701318651ed288dc0368303e5.doc', '2014-09-26 14:50:57.362+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('54257d5199c37bca37000004', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', 'a670c2784a06642141cc37b419bd2b64.html', 'a.html', 586, 'html', 'files/52d26b4e99c37b609a000001/attachs/a670c2784a06642141cc37b419bd2b64.html', '2014-09-26 14:50:57.363+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('54295dae99c37bb079000001', '542674ede5276765ea000000', '52d26b4e99c37b609a000001', '90a1372cd885ea557e02c0d1e6039218.jpeg', 'u=181797442,2321321118&fm=56.jpeg', 16437, 'jpeg', 'files/52d26b4e99c37b609a000001/attachs/90a1372cd885ea557e02c0d1e6039218.jpeg', '2014-09-29 13:25:02.82+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('54295db399c37bb079000002', '542674ede5276765ea000000', '52d26b4e99c37b609a000001', 'f9910b3327ffd6559248809d5bd73cde.jpg', 'u=446392243,3046196010&fm=23&gp=0.jpg', 19382, 'jpg', 'files/52d26b4e99c37b609a000001/attachs/f9910b3327ffd6559248809d5bd73cde.jpg', '2014-09-29 13:25:07.161+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('542a4d4e99c37b20c6000002', '5425390b176ddd41eb000000', '52d8d30799c37b1fb2000001', 'ab2ac8da5a812cfa0fa501f846d73d80.png', 'gemsnote_white.png', 6941, 'png', 'files/52d8d30799c37b1fb2000001/attachs/ab2ac8da5a812cfa0fa501f846d73d80.png', '2014-09-30 06:27:26.236+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('542a643a99c37b5f7a000001', '541ec60e19807a7f65000000', '52d8d30799c37b1fb2000001', '35c7d2cbc4d280b61061efafe310cd14.jpg', '1412007291437.jpg', 2174422, 'jpg', 'files/52d8d30799c37b1fb2000001/attachs/35c7d2cbc4d280b61061efafe310cd14.jpg', '2014-09-30 08:05:14.574+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('542aa1e399c37b7a4f000003', '542a1c59e527674305000000', '52d26b4e99c37b609a000001', '1aedc4b3a944346492dff781b19e7c54.png', 'gemsnote_white.png', 6941, 'png', 'files/52d26b4e99c37b609a000001/attachs/1aedc4b3a944346492dff781b19e7c54.png', '2014-09-30 12:28:19.415+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('542aa28d99c37b7a4f000004', '542a1c59e527674305000000', '52d26b4e99c37b609a000001', '23ef6500d6cad85cbf1d2beb482940be.jpg', 'gemsnote-icon-github副本.jpg', 25570, 'jpg', 'files/52d26b4e99c37b609a000001/attachs/23ef6500d6cad85cbf1d2beb482940be.jpg', '2014-09-30 12:31:09.951+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('542be8e899c37b350e000001', '542bc9b52a30c5120d000001', '52d8d30799c37b1fb2000001', 'ac0ba9ec742c8c50e32d78bda5eac953.png', 'gemsnote_white.png', 6941, 'png', 'files/52d8d30799c37b1fb2000001/attachs/ac0ba9ec742c8c50e32d78bda5eac953.png', '2014-10-01 11:43:36.263+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('542be9ed99c37b350e000002', '542be1d419807a1b50000001', '52d8d30799c37b1fb2000001', 'cd36991efa1c09d4e2832652c9f9ff5f.jpg', 'quick_search_widget_evening.jpg', 10465, 'jpg', 'files/52d8d30799c37b1fb2000001/attachs/cd36991efa1c09d4e2832652c9f9ff5f.jpg', '2014-10-01 11:47:57.686+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('542f7cb299c37b072c000001', '54117af919807a2b8a000003', '52d8d30799c37b1fb2000001', 'b47ad3473db0bb5d9bebe00d9e074e78.log', 'derby.log', 1996, 'log', 'files/52d8d30799c37b1fb2000001/attachs/b47ad3473db0bb5d9bebe00d9e074e78.log', '2014-10-04 04:50:58.404+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('542f7d2f99c37b072c000002', '54117af919807a2b8a000003', '52d8d30799c37b1fb2000001', 'ad732649d48f3f8169bfb1e97c326749.jpg', 'gemsnote-icon-github副本.jpg', 25570, 'jpg', 'files/52d8d30799c37b1fb2000001/attachs/ad732649d48f3f8169bfb1e97c326749.jpg', '2014-10-04 04:53:03.732+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('542f7d3a99c37b072c000003', '54117af919807a2b8a000003', '52d8d30799c37b1fb2000001', 'f9ac80617db98df56ac8280ede4e565b.png', 'gemsnote-64.png', 5799, 'png', 'files/52d8d30799c37b1fb2000001/attachs/f9ac80617db98df56ac8280ede4e565b.png', '2014-10-04 04:53:14.023+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('542f7d3a99c37b072c000004', '54117af919807a2b8a000003', '52d8d30799c37b1fb2000001', '0e2a63826a8c009fef4d06466c58c02d.png', 'gemsnote-icon-github.png', 10930, 'png', 'files/52d8d30799c37b1fb2000001/attachs/0e2a63826a8c009fef4d06466c58c02d.png', '2014-10-04 04:53:14.027+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('542f7d3a99c37b072c000005', '54117af919807a2b8a000003', '52d8d30799c37b1fb2000001', '50131029210c6eb9858be7d0b76c23df.jpg', 'gemsnote-icon-github副本.jpg', 25570, 'jpg', 'files/52d8d30799c37b1fb2000001/attachs/50131029210c6eb9858be7d0b76c23df.jpg', '2014-10-04 04:53:14.031+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('542f7d3a99c37b072c000006', '54117af919807a2b8a000003', '52d8d30799c37b1fb2000001', 'ff4e2348e878205b2018d2b792487456.jpg', 'gemsnote-icon.jpg', 25308, 'jpg', 'files/52d8d30799c37b1fb2000001/attachs/ff4e2348e878205b2018d2b792487456.jpg', '2014-10-04 04:53:14.039+00') ON CONFLICT DO NOTHING;
INSERT INTO public.attachs (id, note_id, upload_user_id, name, title, size, type, path, created_time) VALUES ('542f7d3a99c37b072c000007', '54117af919807a2b8a000003', '52d8d30799c37b1fb2000001', 'f886d97d50ae5d87b66dbb741048aefa.psd', 'gemsnote-icon.psd', 46673, 'psd', 'files/52d8d30799c37b1fb2000001/attachs/f886d97d50ae5d87b66dbb741048aefa.psd', '2014-10-04 04:53:14.046+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: blog_comments; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('542813de99c37b8ff1000002', '5425390b176ddd41eb000000', '52d26b4e99c37b609a000001', '我记得以前可以直接在知乎上看知乎日报，现在怎么不行了？', NULL, '52d8d30799c37b1fb2000001', 0, '{}', '2014-09-28 13:57:50.454+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5428219499c37ba10f000001', '5425390b176ddd41eb000000', '52d26b4e99c37b609a000001', 'ddddddddd', NULL, '52d8d30799c37b1fb2000001', 1, '{52d8d30799c37b1fb2000001}', '2014-09-28 14:56:20.068+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('542821e599c37ba10f000002', '5425390b176ddd41eb000000', '52d26b4e99c37b609a000001', '?????????', NULL, '52d8d30799c37b1fb2000001', 1, '{52d8d30799c37b1fb2000001}', '2014-09-28 14:57:41.711+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('542825b699c37ba73e000002', '5425390b176ddd41eb000000', '52d8d30799c37b1fb2000001', '很好!----------', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-28 15:13:58.413+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('542825c699c37ba73e000003', '5425390b176ddd41eb000000', '52d8d30799c37b1fb2000001', '你怎么这么大家!!!!!!!!!!', NULL, '52d8d30799c37b1fb2000001', 1, '{}', '2014-09-28 15:14:14.809+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('542825c799c37ba73e000004', '5425390b176ddd41eb000000', '52d8d30799c37b1fb2000001', '你怎么这么大家!!!!!!!!!!', NULL, '52d8d30799c37b1fb2000001', 1, '{52d26b4e99c37b609a000001}', '2014-09-28 15:14:15.346+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5428261499c37ba73e000006', '5425390b176ddd41eb000000', '52d3e8ac99c37b7f0d000001', 'xxxxxxxxxxxxx', NULL, '52d8d30799c37b1fb2000001', 2, '{52d8d30799c37b1fb2000001,52d26b4e99c37b609a000001}', '2014-09-28 15:15:32.441+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5428263a99c37ba73e000009', '5425390b176ddd41eb000000', '52d3e8ac99c37b7f0d000001', '好个屁啊!!!!', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-28 15:16:10.421+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54282a1899c37babc8000004', '5425390b176ddd41eb000000', '52d26b4e99c37b609a000001', 'ddd', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-28 15:32:40.621+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5428cd7099c37b1e9300000c', '5342704cdfeb2c27e3000000', '52d8d30799c37b1fb2000001', '我不知道怎么回事?', NULL, '52d26b4e99c37b609a000001', 1, '{52d26b4e99c37b609a000001}', '2014-09-29 03:09:36.05+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5428d18699c37b218b000003', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', '太好了', NULL, '52d8d30799c37b1fb2000001', 0, '{}', '2014-09-29 03:27:02.351+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5428d1ec99c37b218b000005', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', 'xxxxxxxxxxxx', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 03:28:44.679+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5428d1f399c37b218b000006', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', '你好!!', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 03:28:51.685+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5428d41d99c37b218b000011', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', '赞  举报
 赞  举报

 赞  举报', NULL, '52d8d30799c37b1fb2000001', 1, '{52d8d30799c37b1fb2000001}', '2014-09-29 03:38:05.884+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5428d46599c37b218b000013', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', '你懂的我不知道怎么回事?', NULL, '52d8d30799c37b1fb2000001', 1, '{52d8d30799c37b1fb2000001}', '2014-09-29 03:39:17.995+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5428d8dc99c37b218b000023', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', 'xxxxxxxxxxxx', NULL, '52d26b4e99c37b609a000001', 1, '{52d8d30799c37b1fb2000001}', '2014-09-29 03:58:20.946+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5428d92099c37b218b000026', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', 'demo
我不知道怎么回事?
49分钟前 1 赞 删除  回复  取消赞', NULL, '52d8d30799c37b1fb2000001', 1, '{52d8d30799c37b1fb2000001}', '2014-09-29 03:59:28.067+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5428d92c99c37b218b000027', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', '我不知道', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 03:59:40.935+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5428d97999c37b218b000029', '5342704cdfeb2c27e3000000', '52d8d30799c37b1fb2000001', '啦啦', NULL, '52d26b4e99c37b609a000001', 1, '{52d26b4e99c37b609a000001}', '2014-09-29 04:00:57.699+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5428d99399c37b218b00002a', '5342704cdfeb2c27e3000000', '52d8d30799c37b1fb2000001', '人', NULL, '52d26b4e99c37b609a000001', 1, '{52d26b4e99c37b609a000001}', '2014-09-29 04:01:23.624+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5428d9e399c37b218b00002b', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', '专业点赞师!', NULL, '52d8d30799c37b1fb2000001', 0, '{}', '2014-09-29 04:02:43.479+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5428da1399c37b218b00002c', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', 'gemsnote的官方博客--
不一样的笔记! 不一样的我! - gemsnote的官方博客
主页
关于我
博客设置
我的笔记', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 04:03:31.714+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54290c4599c37b5735000004', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', 'xxxxxxxxxxxx', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 07:37:41.551+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54290c4799c37b5735000005', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', 'xxxxxxxxxxxx', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 07:37:43.299+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54290c4899c37b5735000006', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', 'xxxxxxxxxxxx', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 07:37:44.6+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54290c4c99c37b5735000007', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', 'xxxxxxxxxxxx', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 07:37:48.598+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54290c4e99c37b5735000008', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', 'xxxxxxxxxxxx', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 07:37:50.745+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54290c5199c37b5735000009', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 07:37:53.298+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54290c5499c37b573500000a', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 07:37:56.591+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54290c5999c37b573500000b', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 07:38:01.99+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54290c6299c37b573500000c', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 07:38:10.048+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54290c7599c37b573500000d', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', '记录按时间降序排列, 最近的在最前面. 内容只列出了一部分, 点击"展开"即可显示全部; 点击"还原"以该版本还原笔记.', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 07:38:29.2+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54290c7999c37b573500000e', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', '记录按时间降序排列, 最近的在最前面. 内容只列出了一部分, 点击"展开"即可显示全部; 点击"还原"以该版本还原笔记.
记录按时间降序排列, 最近的在最前面. 内容只列出了一部分, 点击"展开"即可显示全部; 点击"还原"以该版本还原笔记.', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 07:38:33.205+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54290c7d99c37b573500000f', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', '记录按时间降序排列, 最近的在最前面. 内容只列出了一部分, 点击"展开"即可显示全部; 点击"还原"以该版本还原笔记.
记录按时间降序排列, 最近的在最前面. 内容只列出了一部分, 点击"展开"即可显示全部; 点击"还原"以该版本还原笔记.', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 07:38:37.534+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('543e150c99c37b3f6c000001', '52d75655a3cf231ad1000000', '52d26b4e99c37b609a000001', 'adfadfads', NULL, NULL, 0, '{}', '2014-10-15 06:32:44.889+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('543e209399c37b4b88000001', '52d75655a3cf231ad1000000', '52d26b4e99c37b609a000001', 'xxxxxxxxxxxx', NULL, NULL, 0, '{}', '2014-10-15 07:21:55.63+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54290c8199c37b5735000010', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', '记录按时间降序排列, 最近的在最前面. 内容只列出了一部分, 点击"展开"即可显示全部; 点击"还原"以该版本还原笔记.
记录按时间降序排列, 最近的在最前面. 内容只列出了一部分, 点击"展开"即可显示全部; 点击"还原"以该版本还原笔记.', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 07:38:41.421+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5429140699c37b6083000004', '54267aa3e527673016000000', '52d26b4e99c37b609a000001', '不错的文章!', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 08:10:46.477+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5429140c99c37b6083000005', '54267aa3e527673016000000', '52d26b4e99c37b609a000001', 'xxxxxxxx', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 08:10:52.697+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5429165599c37b6083000006', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', 'xxxxxxxx', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 08:20:37.469+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5429165a99c37b6083000007', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', 'xxxxxxx', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 08:20:42.331+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54291d7299c37b6bf4000007', '54267aa3e527673016000000', '54291d1999c37b6bf4000002', 'xxxxxxxxxx', NULL, '52d26b4e99c37b609a000001', 1, '{52d26b4e99c37b609a000001}', '2014-09-29 08:50:58.623+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54292a9799c37b7685000004', '542674ede5276765ea000000', '52d8d30799c37b1fb2000001', '你太牛叉了!!', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 09:47:03.145+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54292ae099c37b7685000005', '542674ede5276765ea000000', '52d8d30799c37b1fb2000001', 'xxxxxxxxxx', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 09:48:16.234+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54292c4b99c37b7be8000002', '542674ede5276765ea000000', '52d26b4e99c37b609a000001', 'xxxxxxxxxxxxxxx', NULL, NULL, 0, '{}', '2014-09-29 09:54:19.872+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5429357899c37b818f000001', '542674ede5276765ea000000', '52d26b4e99c37b609a000001', '你好, 欢迎!', NULL, '52d8d30799c37b1fb2000001', 0, '{}', '2014-09-29 10:33:28.331+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54293c7599c37b818f000002', '54267aa3e527673016000000', '52d26b4e99c37b609a000001', 'adadfad', NULL, NULL, 0, '{}', '2014-09-29 11:03:17.711+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('542943ad99c37b818f000003', '531b15b3dfeb2c0ea9000001', '52d26b4e99c37b609a000001', 'xxxxxxxx', NULL, NULL, 0, '{}', '2014-09-29 11:34:05.492+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5429471999c37b818f000005', '54267aa3e527673016000000', '52d26b4e99c37b609a000001', '你不应该这样的!', NULL, '54291d1999c37b6bf4000002', 0, '{}', '2014-09-29 11:48:41.982+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('542948eb99c37b9942000003', '54267aa3e527673016000000', '52d26b4e99c37b609a000001', 'xxxxxxxxxxx', NULL, '54291d1999c37b6bf4000002', 0, '{}', '2014-09-29 11:56:27.188+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54294d5299c37ba03a000001', '542674ede5276765ea000000', '52d8d30799c37b1fb2000001', '这是个坑!!!!', NULL, NULL, 0, '{}', '2014-09-29 12:15:14.076+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54294db399c37ba10a000001', '542674ede5276765ea000000', '52d8d30799c37b1fb2000001', '难道不是吗?', NULL, NULL, 0, '{}', '2014-09-29 12:16:51.427+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54294ecd99c37ba27d000001', '542674ede5276765ea000000', '52d3e8ac99c37b7f0d000001', '信息量：asdfasdfad
kdjfkl

信息量是对信息的度量，就跟温度的度量是摄氏度一样，信息的大小跟随机事件的概率有关。
例如： 在哈尔滨的冬天，一条消息说：哈尔滨明天温度30摄氏度，这个事件肯定会引起轰动，因为它发生的概率很小（信息量大）。日过是夏天，“明天温度30摄氏度”可能没有人觉得是一个新闻，因为夏天温度30摄氏度太正常了，概率太大了（信息点太小了）

从这个例子中可以看出 一个随机事件的信息量的大小与其发生概率是成反相关的。
香农定义的一个事件的信息信息量为：I(X) = log2(1/p) 其中p为事件X发生的概率, 概率越大, 信息量越小', NULL, NULL, 1, '{52d8d30799c37b1fb2000001}', '2014-09-29 12:21:33.182+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54294f1799c37ba27d000002', '542674ede5276765ea000000', '52d26b4e99c37b609a000001', 'PHP登陆后跳转到登陆前页面实现思路及代码_php技巧_脚本之家
PHP登陆后跳转到登陆前页面,利用$_SERVER全局变量可以实现这个功能,下面有个不错的示例,希望对大家有所帮助
www.jb51.net/article/4... 2014-08-08  - 百度快照 - 80%好评
登录后跳转到登录前的页面应该如何做啊_百度知道
1个回答 - 提问时间: 2013年10月12日
最佳答案: 将要跳转的页面以参数的方式传递到登录页面,登录成功后再从参数中提取目标页面 ------解决方案-------------------------------------...
zhidao.baidu.com/link?... 2013-10-13  - 百度快照 - 86%好评', NULL, '52d3e8ac99c37b7f0d000001', 1, '{52d8d30799c37b1fb2000001}', '2014-09-29 12:22:47.896+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54294fab99c37ba3d4000001', '542674ede5276765ea000000', '52d26b4e99c37b609a000001', '用户名或Email
密码 code.google.com/p/go-sqlite/go1/sqlite3
code.google.com/p/go-uuid/uuid
code.google.com/p/go.crypto/bcrypt
code.google.com/p/go.crypto/ssh
code.google.com/p/go.net/context
code.google.com/p/go.net/html
code.google.com/p/go.net/websocket
code.google.com/p/go.tools/present
忘记密码?', NULL, '52d3e8ac99c37b7f0d000001', 0, '{}', '2014-09-29 12:25:15.61+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54294feb99c37ba3d4000002', '542674ede5276765ea000000', '52d8d30799c37b1fb2000001', '你太厉害, 我很佩服!', NULL, '52d3e8ac99c37b7f0d000001', 0, '{}', '2014-09-29 12:26:19.734+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5429503899c37ba4cb000001', '542674ede5276765ea000000', '52d8d30799c37b1fb2000001', 'demo 回复 admin
你太厉害, 我很佩服!
刚刚', NULL, '52d3e8ac99c37b7f0d000001', 0, '{}', '2014-09-29 12:27:36.036+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5429504799c37ba4cb000002', '542674ede5276765ea000000', '52d8d30799c37b1fb2000001', 'gemsnote (作者) 回复 demo
你好, 欢迎!
刚刚  回复  赞  举报', NULL, '52d26b4e99c37b609a000001', 0, '{}', '2014-09-29 12:27:51.583+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5429513b99c37ba673000001', '542674ede5276765ea000000', '52d8d30799c37b1fb2000001', 'xxxxxxxxxx', NULL, NULL, 0, '{}', '2014-09-29 12:31:55.755+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5429515b99c37ba673000002', '542674ede5276765ea000000', '52d26b4e99c37b609a000001', 'xxxxxxxxxxxxx', NULL, '52d8d30799c37b1fb2000001', 0, '{}', '2014-09-29 12:32:27.25+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5429516799c37ba673000003', '542674ede5276765ea000000', '52d26b4e99c37b609a000001', 'xxxxxxxxxxxxxddd', NULL, '52d3e8ac99c37b7f0d000001', 0, '{}', '2014-09-29 12:32:39.235+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5429519699c37ba673000004', '542674ede5276765ea000000', '52d8d30799c37b1fb2000001', '一个新闻，因为夏天温度30摄氏度太正常了，概率太大了（信息点太小了） 从这个例子中可以990', NULL, '52d3e8ac99c37b7f0d000001', 0, '{}', '2014-09-29 12:33:26.032+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('542967e899c37bc41d000002', '542674ede5276765ea000000', '542966f499c37bc034000003', '你懂的', NULL, NULL, 0, '{}', '2014-09-29 14:08:40.031+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('54353f4f99c37b59bc000001', '54352ffe99c37b423c000005', '52d26b4e99c37b609a000001', 'ddddddddddddddd', NULL, NULL, 0, '{}', '2014-10-08 13:42:39.925+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5435402999c37b59bc000002', '54352ffe99c37b423c000005', '52d26b4e99c37b609a000001', '欢迎来到gemsnote!
快速链接

我的笔记
gemsnote 登录
gemsnote 主页
lea++, gemsnote博客平台
gemsnote 社区
gemsnote github', NULL, NULL, 0, '{}', '2014-10-08 13:46:17.55+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5435407399c37b5b71000001', '54352ffe99c37b423c000005', '52d26b4e99c37b609a000001', '+, gemsnote博客平台 gemsnote 社区 gemsnote github
刚刚  删除', NULL, NULL, 0, '{}', '2014-10-08 13:47:31.48+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5438a62f99c37b2ab6000003', '541eb10019807a3063000000', '52d8d30799c37b1fb2000001', 'kjlk, you can make it!', NULL, NULL, 0, '{}', '2014-10-11 03:38:23.745+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('543a852499c37b14d7000001', '541eb0e219807a0bfc000000', '52d8d30799c37b1fb2000001', 'Hello, you can make it!!', NULL, NULL, 0, '{}', '2014-10-12 13:41:56.223+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('543a9a9e99c37b3bcf000001', '541eb0e219807a0bfc000000', '52d8d30799c37b1fb2000001', 'xxxxxxxxxxxx', NULL, NULL, 0, '{}', '2014-10-12 15:13:34.428+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_comments (id, note_id, user_id, content, to_comment_id, to_user_id, like_num, like_user_ids, created_time) VALUES ('5478168ac1cc923150000002', '5368c1b919807a6f95000000', '5368c1aa99c37b029d000001', 'dfsd', NULL, NULL, 0, '{}', '2014-11-28 06:30:34.082+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: blog_likes; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.blog_likes (id, note_id, user_id, created_time) VALUES ('542825e099c37ba73e000005', '5425390b176ddd41eb000000', '52d8d30799c37b1fb2000001', '2014-09-28 15:14:40.832+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_likes (id, note_id, user_id, created_time) VALUES ('5428265699c37ba73e00000c', '5425390b176ddd41eb000000', '52d3e8ac99c37b7f0d000001', '2014-09-28 15:16:38.897+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_likes (id, note_id, user_id, created_time) VALUES ('5428432499c37bcabb000001', '5425390b176ddd41eb000000', '52d26b4e99c37b609a000001', '2014-09-28 17:19:32.861+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_likes (id, note_id, user_id, created_time) VALUES ('5428cc0399c37b1e93000002', '5342704cdfeb2c27e3000000', '52d8d30799c37b1fb2000001', '2014-09-29 03:03:31.882+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_likes (id, note_id, user_id, created_time) VALUES ('5428d64399c37b218b00001c', '5342704cdfeb2c27e3000000', '52d26b4e99c37b609a000001', '2014-09-29 03:47:15.649+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_likes (id, note_id, user_id, created_time) VALUES ('54291d9f99c37b6bf4000009', '54267aa3e527673016000000', '52d8d30799c37b1fb2000001', '2014-09-29 08:51:43.665+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_likes (id, note_id, user_id, created_time) VALUES ('542946f799c37b818f000004', '54267aa3e527673016000000', '52d26b4e99c37b609a000001', '2014-09-29 11:48:07.823+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_likes (id, note_id, user_id, created_time) VALUES ('542a121599c37b0f88000003', '531b15b3dfeb2c0ea9000001', '52d26b4e99c37b609a000001', '2014-09-30 02:14:45.08+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_likes (id, note_id, user_id, created_time) VALUES ('542a131199c37b0f88000004', '531b15b3dfeb2c0ea9000001', '52d8d30799c37b1fb2000001', '2014-09-30 02:18:57.119+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_likes (id, note_id, user_id, created_time) VALUES ('5438a61299c37b2ab6000002', '541eb10019807a3063000000', '52d8d30799c37b1fb2000001', '2014-10-11 03:37:54.28+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_likes (id, note_id, user_id, created_time) VALUES ('543a9c4b99c37b3bcf000002', '541eb0e219807a0bfc000000', '52d8d30799c37b1fb2000001', '2014-10-12 15:20:43.124+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_likes (id, note_id, user_id, created_time) VALUES ('543cd1c099c37b535c000001', '543ccd6de6bddb06d100000d', '52d26b4e99c37b609a000001', '2014-10-14 07:33:20.417+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_likes (id, note_id, user_id, created_time) VALUES ('543e20a099c37b4b88000002', '52d75655a3cf231ad1000000', '52d26b4e99c37b609a000001', '2014-10-15 07:22:08.919+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: blog_singles; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.blog_singles (id, user_id, title, url_title, content, updated_time, created_time) VALUES ('546325e699c37b80ae000001', '5368c1aa99c37b029d000001', 'About Me', 'About-Me', '<p>Hello,&nbsp;I am Gemsnote (^_^).</p>', '2015-06-15 10:38:51.668+00', '2014-11-12 09:18:30.046+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_singles (id, user_id, title, url_title, content, updated_time, created_time) VALUES ('546325e699c37b80ae000002', '5368c9fc99c37b095a000006', 'About Me', 'About-Me', '<p>Hello, I am (^_^)</p>', '2014-11-12 09:18:30.051+00', '2014-11-12 09:18:30.051+00') ON CONFLICT DO NOTHING;
INSERT INTO public.blog_singles (id, user_id, title, url_title, content, updated_time, created_time) VALUES ('5524ba2f99c37b292000000c', '5524ba2f99c37b2920000007', 'About Me', 'About-Me', 'Hello, I am (^_^)', '2015-04-08 05:18:39.444+00', '2015-04-08 05:18:39.444+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: blogs; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: configs; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('52d26b4e99c37b609a000001', '5368c1aa99c37b029d000001', '', '', '{}', '{}', '[]', false, false, false, '0001-01-01 00:00:00+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543685ff99c37b7fe7000001', '5368c1aa99c37b029d000001', 'toImageBinPath', 'lllllllllll', '{}', '{}', '[]', false, false, false, '2014-10-09 12:56:31.894+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436861399c37b7fe7000002', '5368c1aa99c37b029d000001', 'noteSubDomain', '', '{}', '{}', '[]', false, false, false, '2014-10-20 12:03:22.347+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436861399c37b7fe7000003', '5368c1aa99c37b029d000001', 'blogSubDomain', '', '{}', '{}', '[]', false, false, false, '2014-10-20 12:03:22.36+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436861399c37b7fe7000004', '5368c1aa99c37b029d000001', 'leaSubDomain', '', '{}', '{}', '[]', false, false, false, '2014-10-20 12:03:22.36+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543686aa99c37b80d0000001', '5368c1aa99c37b029d000001', 'recommendTags', '', '{小写,golang,gemsnote}', '{}', '[]', true, false, false, '2014-10-09 12:59:22.155+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543686aa99c37b80d0000002', '5368c1aa99c37b029d000001', 'newTags', '', '{小写,golang,gemsnote,haha}', '{}', '[]', true, false, false, '2014-10-09 12:59:22.155+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543686e399c37b810f000001', '5368c1aa99c37b029d000001', 'emailHost', 'smtp.ym.163.com', '{}', '{}', '[]', false, false, false, '2014-10-22 12:13:08.635+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543686e399c37b810f000002', '5368c1aa99c37b029d000001', 'emailPort', '25', '{}', '{}', '[]', false, false, false, '2014-10-22 12:13:08.636+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543686e399c37b810f000003', '5368c1aa99c37b029d000001', 'emailUsername', '', '{}', '{}', '[]', false, false, false, '2014-10-22 12:13:08.636+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543686e399c37b810f000004', '5368c1aa99c37b029d000001', 'emailPassword', '', '{}', '{}', '[]', false, false, false, '2014-10-22 12:13:08.636+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687a799c37b81f6000001', '5368c1aa99c37b029d000001', 'emailTemplateHeader', '<div style="width: 600px; margin:auto; border-radius:5px; border: 1px solid #ccc; padding: 20px;">
			<div>
				<div>
					<div style="float:left; height: 40px;">
						<a href="{{$.siteUrl}}" style="font-size: 24px">gemsnote</a>
					</div>
					<div style="float:left; height:40px; line-height:40px;">
						&nbsp;&nbsp;| &nbsp;<span style="font-size:14px">{{$.subject}}</span>
					</div>
					<div style="clear:both"></div>
				</div>
			</div>
			<hr style="border:none;border-top: 1px solid #ccc"/>
			<div style="margin-top: 20px; font-size: 14px;">
				', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.149+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687a799c37b81f6000002', '5368c1aa99c37b029d000001', 'emailTemplateFooter', '</div>

			<div id="gemsnoteFooter" style="margin-top: 30px; border-top: 1px solid #ccc">
				<style>
					#gemsnoteFooter {
						color: #666;
						font-size: 12px;
					}
					#gemsnoteFooter a {
						color: #666;
						font-size: 12px;
					}
				</style>
				<a href="{{$.siteUrl}}">gemsnote</a>, your own cloud note!
			</div>
		</div>', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.15+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687a799c37b81f6000003', '5368c1aa99c37b029d000001', 'emailTemplateRegisterSubject', '欢迎来到gemsnote, 请验证邮箱', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.15+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687a799c37b81f6000004', '5368c1aa99c37b029d000001', 'emailTemplateRegister', '{{header}}
<p>
{{$.user.email}} 您好, 欢迎来到gemsnote.
</p>
<p>
请点击链接验证邮箱: <a href="{{$.tokenUrl}}">{{$.tokenUrl}}</a>
</p>
<p>
{{$.tokenTimeout}}小时后过期.
</p>
{{footer}}', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.151+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687a799c37b81f6000005', '5368c1aa99c37b029d000001', 'emailTemplateFindPasswordSubject', '找回密码', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.151+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687a799c37b81f6000006', '5368c1aa99c37b029d000001', 'emailTemplateFindPassword', '{{header}}
<p>
请点击链接修改密码 <a href="{{$.tokenUrl}}">{{$.tokenUrl}}</a>
</p>
<p>
{{$.tokenTimeout}}小时后过期.
</p>

{{footer}}', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.152+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687a799c37b81f6000007', '5368c1aa99c37b029d000001', 'emailTemplateUpdateEmailSubject', '验证邮箱', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.152+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687a799c37b81f6000008', '5368c1aa99c37b029d000001', 'emailTemplateUpdateEmail', '{{header}}
<p>
邮箱验证后您的登录邮箱为: {{$.newEmail}}
</p>
<p>
请点击链接验证邮箱: <a href="{{$.tokenUrl}}">{{$.tokenUrl}}</a>
</p>
<p>
{{$.tokenTimeout}}小时后过期.
</p>
{{footer}}
', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.152+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687a799c37b81f6000009', '5368c1aa99c37b029d000001', 'emailTemplateInviteSubject', '邀请注册gemsnote', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.152+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687a799c37b81f600000a', '5368c1aa99c37b029d000001', 'emailTemplateInvite', '{{header}}

<p>您好, 您的好友{{$.user.email}}邀请您注册gemsnote</p>

<p>Ta的留言: {{$.content}}</p>

<p>点击链接注册gemsnote <a href="{{$.registerUrl}}">{{$.registerUrl}}</a></p>

{{footer}}
', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.153+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687a799c37b81f600000b', '5368c1aa99c37b029d000001', 'emailTemplateCommentSubject', '评论提醒', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.153+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687a799c37b81f600000c', '5368c1aa99c37b029d000001', 'emailTemplateComment', '{{header}}
<p>
{{if $.commentedUser.isBlogAuthor}}
您的博客 "{{$.blog.title}}" 被 {{$.commentUser.username}} 评论了.
{{else}}
您在 "{{$.blog.title}}" 发表的评论被 {{$.commentUser.username}}{{if $.commentUser.isBlogAuthor}}(作者){{end}} 评论了.
{{end}}
</p>

<div>
<b>评论内容: </b>
<blockquote>{{$.commentContent}}</blockquote>
</div>
<p>
博客链接: <a href="{{$.blog.url}}">{{$.blog.url}}</a>
</p>
{{footer}} ', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.153+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687c599c37b81f600000d', '5368c1aa99c37b029d000001', 'emailTemplateHeader', '<div style="width: 600px; margin:auto; border-radius:5px; border: 1px solid #ccc; padding: 20px;">
			<div>
				<div>
					<div style="float:left; height: 40px;">
						<a href="{{$.siteUrl}}" style="font-size: 24px">gemsnote</a>
					</div>
					<div style="float:left; height:40px; line-height:40px;">
						&nbsp;&nbsp;| &nbsp;<span style="font-size:14px">{{$.subject}}</span>
					</div>
					<div style="clear:both"></div>
				</div>
			</div>
			<hr style="border:none;border-top: 1px solid #ccc"/>
			<div style="margin-top: 20px; font-size: 14px;">
				', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.149+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687c599c37b81f600000e', '5368c1aa99c37b029d000001', 'emailTemplateFooter', '</div>

			<div id="gemsnoteFooter" style="margin-top: 30px; border-top: 1px solid #ccc">
				<style>
					#gemsnoteFooter {
						color: #666;
						font-size: 12px;
					}
					#gemsnoteFooter a {
						color: #666;
						font-size: 12px;
					}
				</style>
				<a href="{{$.siteUrl}}">gemsnote</a>, your own cloud note!
			</div>
		</div>', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.15+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687c599c37b81f600000f', '5368c1aa99c37b029d000001', 'emailTemplateRegisterSubject', '欢迎来到gemsnote, 请验证邮箱', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.15+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687c599c37b81f6000010', '5368c1aa99c37b029d000001', 'emailTemplateRegister', '{{header}}
<p>
{{$.user.email}} 您好, 欢迎来到gemsnote.
</p>
<p>
请点击链接验证邮箱: <a href="{{$.tokenUrl}}">{{$.tokenUrl}}</a>
</p>
<p>
{{$.tokenTimeout}}小时后过期.
</p>
{{footer}}', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.151+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687c599c37b81f6000011', '5368c1aa99c37b029d000001', 'emailTemplateFindPasswordSubject', '找回密码', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.151+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687c599c37b81f6000012', '5368c1aa99c37b029d000001', 'emailTemplateFindPassword', '{{header}}
<p>
请点击链接修改密码 <a href="{{$.tokenUrl}}">{{$.tokenUrl}}</a>
</p>
<p>
{{$.tokenTimeout}}小时后过期.
</p>

{{footer}}', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.152+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687c599c37b81f6000013', '5368c1aa99c37b029d000001', 'emailTemplateUpdateEmailSubject', '验证邮箱', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.152+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('546326d499c37b80ae000015', '5368c1aa99c37b029d000001', 'uploadAvatarSize', '1', '{}', '{}', '[]', false, false, false, '2014-11-12 09:22:28.397+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687c599c37b81f6000014', '5368c1aa99c37b029d000001', 'emailTemplateUpdateEmail', '{{header}}
<p>
邮箱验证后您的登录邮箱为: {{$.newEmail}}
</p>
<p>
请点击链接验证邮箱: <a href="{{$.tokenUrl}}">{{$.tokenUrl}}</a>
</p>
<p>
{{$.tokenTimeout}}小时后过期.
</p>
{{footer}}
', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.152+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687c599c37b81f6000015', '5368c1aa99c37b029d000001', 'emailTemplateInviteSubject', '邀请注册gemsnote', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.152+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687c599c37b81f6000016', '5368c1aa99c37b029d000001', 'emailTemplateInvite', '{{header}}

<p>您好, 您的好友{{$.user.email}}邀请您注册gemsnote</p>

<p>Ta的留言: {{$.content}}</p>

<p>点击链接注册gemsnote <a href="{{$.registerUrl}}">{{$.registerUrl}}</a></p>

{{footer}}
', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.153+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687c599c37b81f6000017', '5368c1aa99c37b029d000001', 'emailTemplateCommentSubject', '评论提醒', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.153+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543687c599c37b81f6000018', '5368c1aa99c37b029d000001', 'emailTemplateComment', '{{header}}
<p>
{{if $.commentedUser.isBlogAuthor}}
您的博客 "{{$.blog.title}}" 被 {{$.commentUser.username}} 评论了.
{{else}}
您在 "{{$.blog.title}}" 发表的评论被 {{$.commentUser.username}}{{if $.commentUser.isBlogAuthor}}(作者){{end}} 评论了.
{{end}}
</p>

<div>
<b>评论内容: </b>
<blockquote>{{$.commentContent}}</blockquote>
</div>
<p>
博客链接: <a href="{{$.blog.url}}">{{$.blog.url}}</a>
</p>
{{footer}} ', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.153+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436882499c37b81f6000019', '5368c1aa99c37b029d000001', 'emailTemplateHeader', '<div style="width: 600px; margin:auto; border-radius:5px; border: 1px solid #ccc; padding: 20px;">
			<div>
				<div>
					<div style="float:left; height: 40px;">
						<a href="{{$.siteUrl}}" style="font-size: 24px">gemsnote</a>
					</div>
					<div style="float:left; height:40px; line-height:40px;">
						&nbsp;&nbsp;| &nbsp;<span style="font-size:14px">{{$.subject}}</span>
					</div>
					<div style="clear:both"></div>
				</div>
			</div>
			<hr style="border:none;border-top: 1px solid #ccc"/>
			<div style="margin-top: 20px; font-size: 14px;">
				', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.149+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436882499c37b81f600001a', '5368c1aa99c37b029d000001', 'emailTemplateFooter', '</div>

			<div id="gemsnoteFooter" style="margin-top: 30px; border-top: 1px solid #ccc">
				<style>
					#gemsnoteFooter {
						color: #666;
						font-size: 12px;
					}
					#gemsnoteFooter a {
						color: #666;
						font-size: 12px;
					}
				</style>
				<a href="{{$.siteUrl}}">gemsnote</a>, your own cloud note!
			</div>
		</div>', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.15+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436882499c37b81f600001b', '5368c1aa99c37b029d000001', 'emailTemplateRegisterSubject', '欢迎来到gemsnote, 请验证邮箱', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.15+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436882499c37b81f600001c', '5368c1aa99c37b029d000001', 'emailTemplateRegister', '{{header}}
<p>
{{$.user.email}} 您好, 欢迎来到gemsnote.
</p>
<p>
请点击链接验证邮箱: <a href="{{$.tokenUrl}}">{{$.tokenUrl}}</a>
</p>
<p>
{{$.tokenTimeout}}小时后过期.
</p>
{{footer}}', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.151+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436882499c37b81f600001d', '5368c1aa99c37b029d000001', 'emailTemplateFindPasswordSubject', '找回密码', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.151+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436882499c37b81f600001e', '5368c1aa99c37b029d000001', 'emailTemplateFindPassword', '{{header}}
<p>
请点击链接修改密码 <a href="{{$.tokenUrl}}">{{$.tokenUrl}}</a>
</p>
<p>
{{$.tokenTimeout}}小时后过期.
</p>

{{footer}}', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.152+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436882499c37b81f600001f', '5368c1aa99c37b029d000001', 'emailTemplateUpdateEmailSubject', '验证邮箱', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.152+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436882499c37b81f6000020', '5368c1aa99c37b029d000001', 'emailTemplateUpdateEmail', '{{header}}
<p>
邮箱验证后您的登录邮箱为: {{$.newEmail}}
</p>
<p>
请点击链接验证邮箱: <a href="{{$.tokenUrl}}">{{$.tokenUrl}}</a>
</p>
<p>
{{$.tokenTimeout}}小时后过期.
</p>
{{footer}}
', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.152+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436882499c37b81f6000021', '5368c1aa99c37b029d000001', 'emailTemplateInviteSubject', '邀请注册gemsnote', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.152+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436882499c37b81f6000022', '5368c1aa99c37b029d000001', 'emailTemplateInvite', '{{header}}

<p>您好, 您的好友{{$.user.email}}邀请您注册gemsnote</p>

<p>Ta的留言: {{$.content}}</p>

<p>点击链接注册gemsnote <a href="{{$.registerUrl}}">{{$.registerUrl}}</a></p>

{{footer}}
', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.153+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436882499c37b81f6000023', '5368c1aa99c37b029d000001', 'emailTemplateCommentSubject', '评论提醒', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.153+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436882499c37b81f6000024', '5368c1aa99c37b029d000001', 'emailTemplateComment', '{{header}}
<p>
{{if $.commentedUser.isBlogAuthor}}
您的博客 "{{$.blog.title}}" 被 {{$.commentUser.username}} 评论了.
{{else}}
您在 "{{$.blog.title}}" 发表的评论被 {{$.commentUser.username}}{{if $.commentUser.isBlogAuthor}}(作者){{end}} 评论了.
{{end}}
</p>

<div>
<b>评论内容: </b>
<blockquote>{{$.commentContent}}</blockquote>
</div>
<p>
博客链接: <a href="{{$.blog.url}}">{{$.blog.url}}</a>
</p>
{{footer}} ', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.153+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436889899c37b81f6000025', '5368c1aa99c37b029d000001', 'emailTemplateHeader', '<div style="width: 600px; margin:auto; border-radius:5px; border: 1px solid #ccc; padding: 20px;">
			<div>
				<div>
					<div style="float:left; height: 40px;">
						<a href="{{$.siteUrl}}" style="font-size: 24px">gemsnote</a>
					</div>
					<div style="float:left; height:40px; line-height:40px;">
						&nbsp;&nbsp;| &nbsp;<span style="font-size:14px">{{$.subject}}</span>
					</div>
					<div style="clear:both"></div>
				</div>
			</div>
			<hr style="border:none;border-top: 1px solid #ccc"/>
			<div style="margin-top: 20px; font-size: 14px;">
				', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.149+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436889899c37b81f6000026', '5368c1aa99c37b029d000001', 'emailTemplateFooter', '</div>

			<div id="gemsnoteFooter" style="margin-top: 30px; border-top: 1px solid #ccc">
				<style>
					#gemsnoteFooter {
						color: #666;
						font-size: 12px;
					}
					#gemsnoteFooter a {
						color: #666;
						font-size: 12px;
					}
				</style>
				<a href="{{$.siteUrl}}">gemsnote</a>, your own cloud note!
			</div>
		</div>', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.15+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436889899c37b81f6000027', '5368c1aa99c37b029d000001', 'emailTemplateRegisterSubject', '欢迎来到gemsnote, 请验证邮箱', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.15+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436889899c37b81f6000028', '5368c1aa99c37b029d000001', 'emailTemplateRegister', '{{header}}
<p>
{{$.user.email}} 您好, 欢迎来到gemsnote.
</p>
<p>
请点击链接验证邮箱: <a href="{{$.tokenUrl}}">{{$.tokenUrl}}</a>
</p>
<p>
{{$.tokenTimeout}}小时后过期.
</p>
{{footer}}', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.151+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436889899c37b81f6000029', '5368c1aa99c37b029d000001', 'emailTemplateFindPasswordSubject', '找回密码', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.151+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436889899c37b81f600002a', '5368c1aa99c37b029d000001', 'emailTemplateFindPassword', '{{header}}
<p>
请点击链接修改密码 <a href="{{$.tokenUrl}}">{{$.tokenUrl}}</a>
</p>
<p>
{{$.tokenTimeout}}小时后过期.
</p>

{{footer}}', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.152+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436889899c37b81f600002b', '5368c1aa99c37b029d000001', 'emailTemplateUpdateEmailSubject', '验证邮箱', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.152+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436889899c37b81f600002c', '5368c1aa99c37b029d000001', 'emailTemplateUpdateEmail', '{{header}}
<p>
邮箱验证后您的登录邮箱为: {{$.newEmail}}
</p>
<p>
请点击链接验证邮箱: <a href="{{$.tokenUrl}}">{{$.tokenUrl}}</a>
</p>
<p>
{{$.tokenTimeout}}小时后过期.
</p>
{{footer}}
', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.152+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436889899c37b81f600002d', '5368c1aa99c37b029d000001', 'emailTemplateInviteSubject', '邀请注册gemsnote', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.152+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436889899c37b81f600002e', '5368c1aa99c37b029d000001', 'emailTemplateInvite', '{{header}}

<p>您好, 您的好友{{$.user.email}}邀请您注册gemsnote</p>

<p>Ta的留言: {{$.content}}</p>

<p>点击链接注册gemsnote <a href="{{$.registerUrl}}">{{$.registerUrl}}</a></p>

{{footer}}
', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.153+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436889899c37b81f600002f', '5368c1aa99c37b029d000001', 'emailTemplateCommentSubject', '评论提醒', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.153+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436889899c37b81f6000030', '5368c1aa99c37b029d000001', 'emailTemplateComment', '{{header}}
<p>
{{if $.commentedUser.isBlogAuthor}}
您的博客 "{{$.blog.title}}" 被 {{$.commentUser.username}} 评论了.
{{else}}
您在 "{{$.blog.title}}" 发表的评论被 {{$.commentUser.username}}{{if $.commentUser.isBlogAuthor}}(作者){{end}} 评论了.
{{end}}
</p>

<div>
<b>评论内容: </b>
<blockquote>{{$.commentContent}}</blockquote>
</div>
<p>
博客链接: <a href="{{$.blog.url}}">{{$.blog.url}}</a>
</p>
{{footer}} ', '{}', '{}', '[]', false, false, false, '2014-10-09 14:06:53.153+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436992c99c37b9b96000001', '5368c1aa99c37b029d000001', 'userFilterEmail', '', '{}', '{}', '[]', false, false, false, '2014-10-22 11:53:11.258+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436992c99c37b9b96000002', '5368c1aa99c37b029d000001', 'userFilterWhiteList', 'lifephp@gmail.com
life@gemsnote.com', '{}', '{}', '[]', false, false, false, '2014-10-22 11:53:11.258+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5436992c99c37b9b96000003', '5368c1aa99c37b029d000001', 'userFilterBlackList', '', '{}', '{}', '[]', false, false, false, '2014-10-22 11:53:11.259+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('54369f2e99c37ba4c0000001', '5368c1aa99c37b029d000001', 'latestEmailSubject', '{{$.username}} 欢迎来到gemsnote', '{}', '{}', '[]', false, false, false, '2014-10-22 11:57:57.581+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('54369f2e99c37ba4c0000002', '5368c1aa99c37b029d000001', 'latestEmailBody', '{{header}}
{{$.username}}
欢迎来到gemsnote!
{{footer}}', '{}', '{}', '[]', false, false, false, '2014-10-22 11:57:57.581+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5437a18c99c37b7198000003', '5368c1aa99c37b029d000001', 'oldEmails', '', '{}', '{"{{$.username}} 欢迎来到gemsnote": "{{header}}\r\n{{$.username}} 欢迎来到gemsnote\r\n{{footer}}"}', '[]', false, true, false, '2014-10-10 09:06:20.723+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5437a68899c37b78cf000001', '5368c1aa99c37b029d000001', 'sendEmails', '', '{}', '{}', '[]', false, false, false, '2014-10-22 11:57:57.58+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5437a6ae99c37b78cf000003', '5368c1aa99c37b029d000001', 'sendEmails', '', '{}', '{}', '[]', false, false, false, '2014-10-22 11:57:57.58+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5437a6e199c37b78cf000006', '5368c1aa99c37b029d000001', 'sendEmails', '', '{}', '{}', '[]', false, false, false, '2014-10-22 11:57:57.58+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5437da5099c37ba810000001', '5368c1aa99c37b029d000001', 'registerSharedUserId', '5368c1aa99c37b029d000001', '{}', '{}', '[]', false, false, false, '2015-04-08 05:18:19.909+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5437da5399c37ba810000002', '5368c1aa99c37b029d000001', 'registerSharedUserId', '5368c1aa99c37b029d000001', '{}', '{}', '[]', false, false, false, '2015-04-08 05:18:19.909+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5437da6c99c37ba810000003', '5368c1aa99c37b029d000001', 'registerSharedUserId', '5368c1aa99c37b029d000001', '{}', '{}', '[]', false, false, false, '2015-04-08 05:18:19.909+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5437da9e99c37ba8b4000001', '5368c1aa99c37b029d000001', 'registerSharedNotebooks', '', '{}', '{}', '[]', false, false, true, '2015-04-08 05:18:19.912+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5437db0699c37ba8b4000003', '5368c1aa99c37b029d000001', 'registerSharedNotebooks', '', '{}', '{}', '[]', false, false, true, '2015-04-08 05:18:19.912+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5437db3399c37ba8b4000005', '5368c1aa99c37b029d000001', 'registerSharedNotebooks', '', '{}', '{}', '[]', false, false, true, '2015-04-08 05:18:19.912+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543ce20999c37b69f7000001', '5368c1aa99c37b029d000001', 'allowCustomDomain', '', '{}', '{}', '[]', false, false, false, '2014-10-20 12:03:22.361+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543ce20a99c37b69f7000002', '5368c1aa99c37b029d000001', 'blackSubDomains', '', '{note,blog,lea,gemsnote,gogo}', '{}', '[]', true, false, false, '2014-10-20 12:03:22.361+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('543ce20a99c37b69f7000003', '5368c1aa99c37b029d000001', 'blackCustomDomains', '', '{gemsnote.com,lealife.com,gemsnote.cn,gemsnote.org,gemsnote.net}', '{}', '[]', true, false, false, '2014-10-20 12:03:22.364+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('54468afc99c37b528b000001', '5368c1aa99c37b029d000001', 'mongodumpPath', '/Users/life/Desktop/hadoop/mongodb-osx-x86_64-2.4.7/bin/mongodump', '{}', '{}', '[]', false, false, false, '2014-11-12 09:23:42.849+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('54468afc99c37b528b000002', '5368c1aa99c37b029d000001', 'mongorestorePath', '/Users/life/Desktop/hadoop/mongodb-osx-x86_64-2.4.7/bin/mongorestore', '{}', '{}', '[]', false, false, false, '2014-11-12 09:23:42.85+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5447098999c37b0ef6000001', '5368c1aa99c37b029d000001', 'backups', '', '{}', '{}', '[]', false, false, true, '2014-11-12 11:35:49.989+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5447099299c37b0ef6000002', '5368c1aa99c37b029d000001', 'backups', '', '{}', '{}', '[]', false, false, true, '2014-11-12 11:35:49.989+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('54485f4699c37b0585000001', '5368c1aa99c37b029d000001', 'openRegister', '1', '{}', '{}', '[]', false, false, false, '2014-10-23 02:01:09.178+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('546325e699c37b80ae000007', '5368c1aa99c37b029d000001', 'UpgradeBetaToBeta2', '1', '{}', '{}', '[]', false, false, false, '2014-11-12 10:49:33.883+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5463269599c37b80ae00000d', '5368c1aa99c37b029d000001', 'homePage', '1', '{}', '{}', '[]', false, false, false, '2015-06-15 10:50:55.279+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('5463269899c37b80ae00000e', '5368c1aa99c37b029d000001', 'homePage', '1', '{}', '{}', '[]', false, false, false, '2015-06-15 10:50:55.279+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('546326a199c37b80ae00000f', '5368c1aa99c37b029d000001', 'homePage', '1', '{}', '{}', '[]', false, false, false, '2015-06-15 10:50:55.279+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('546326b199c37b80ae000010', '5368c1aa99c37b029d000001', 'uploadImageSize', '1', '{}', '{}', '[]', false, false, false, '2014-11-12 09:21:53.157+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('546326b199c37b80ae000011', '5368c1aa99c37b029d000001', 'uploadAvatarSize', '1', '{}', '{}', '[]', false, false, false, '2014-11-12 09:21:53.158+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('546326b199c37b80ae000012', '5368c1aa99c37b029d000001', 'uploadBlogLogoSize', '1', '{}', '{}', '[]', false, false, false, '2014-11-12 09:21:53.159+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('546326b199c37b80ae000013', '5368c1aa99c37b029d000001', 'uploadAttachSize', '1', '{}', '{}', '[]', false, false, false, '2014-11-12 09:21:53.16+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('546326d499c37b80ae000014', '5368c1aa99c37b029d000001', 'uploadImageSize', '0.01', '{}', '{}', '[]', false, false, false, '2014-11-12 09:22:28.396+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('546326d499c37b80ae000016', '5368c1aa99c37b029d000001', 'uploadBlogLogoSize', '1', '{}', '{}', '[]', false, false, false, '2014-11-12 09:22:28.397+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('546326d499c37b80ae000017', '5368c1aa99c37b029d000001', 'uploadAttachSize', '1', '{}', '{}', '[]', false, false, false, '2014-11-12 09:22:28.398+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('546326e799c37b80ae000018', '5368c1aa99c37b029d000001', 'uploadImageSize', '1', '{}', '{}', '[]', false, false, false, '2014-11-12 09:22:47.196+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('546326e799c37b80ae000019', '5368c1aa99c37b029d000001', 'uploadAvatarSize', '1', '{}', '{}', '[]', false, false, false, '2014-11-12 09:22:47.198+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('546326e799c37b80ae00001a', '5368c1aa99c37b029d000001', 'uploadBlogLogoSize', '1', '{}', '{}', '[]', false, false, false, '2014-11-12 09:22:47.198+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('546326e799c37b80ae00001b', '5368c1aa99c37b029d000001', 'uploadAttachSize', '10', '{}', '{}', '[]', false, false, false, '2014-11-12 09:22:47.199+00') ON CONFLICT DO NOTHING;
INSERT INTO public.configs (id, user_id, key, value_str, value_arr, value_map, value_arr_map, is_arr, is_map, is_arr_map, updated_time) VALUES ('551a3f5c99c37b04de000004', '5368c1aa99c37b029d000001', 'UpgradeBetaToBeta4', '1', '{}', '{}', '[]', false, false, false, '2015-03-31 06:31:56.608+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: files; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.files (id, user_id, album_id, name, title, size, type, path, is_default_album, created_time, from_file_id) VALUES ('54081515acf7541eaf8369cb', '5368c9fc99c37b095a000006', '52d3e8ac99c37b7f0d000001', 'f5ed4d84df526bf788f21ce4c11439ca.png', 'gemsnote_green.png', 15799, '', '/upload/5368c9fc99c37b095a000006/images/f5ed4d84df526bf788f21ce4c11439ca.png', true, '2014-09-04 07:30:29.508+00', NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.files (id, user_id, album_id, name, title, size, type, path, is_default_album, created_time, from_file_id) VALUES ('5408152facf7541eaf8369cc', '5368c9fc99c37b095a000006', '52d3e8ac99c37b7f0d000001', 'c93d600ee42c94876258cd0d51b36e06.png', 'gemsnote-icon-github.png', 10930, '', '/upload/5368c9fc99c37b095a000006/images/c93d600ee42c94876258cd0d51b36e06.png', true, '2014-09-04 07:30:55.23+00', NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.files (id, user_id, album_id, name, title, size, type, path, is_default_album, created_time, from_file_id) VALUES ('561bb56899c37b2f1e000001', '5368c1aa99c37b029d000001', '52d3e8ac99c37b7f0d000001', '2e1801587a1ade1b54d4e31f26d07b33.jpg', 'gemsnote-icon.jpg', 25308, '', 'public/upload/517/5368c1aa99c37b029d000001/images/logo/2e1801587a1ade1b54d4e31f26d07b33.jpg', true, '2015-10-12 13:28:08.408+00', NULL) ON CONFLICT DO NOTHING;


--
-- Data for Name: group_users; Type: TABLE DATA; Schema: public;


--
-- Data for Name: groups; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.groups (id, user_id, title, user_count, created_time, users) VALUES ('5463263299c37b80ae000009', '5368c1aa99c37b029d000001', 'love gemsnote', 0, '2014-11-12 09:19:46.664+00', NULL) ON CONFLICT DO NOTHING;


--
-- Data for Name: has_share_notes; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.has_share_notes (id, user_id, to_user_id, seq) VALUES ('5368c9fcdd73244c2b52e10a', '5368c1aa99c37b029d000001', '5368c9fc99c37b095a000006', 0) ON CONFLICT DO NOTHING;
INSERT INTO public.has_share_notes (id, user_id, to_user_id, seq) VALUES ('540814f1acf7541eaf8369c7', '5368c1aa99c37b029d000001', '540814f199c37b555d000001', 0) ON CONFLICT DO NOTHING;
INSERT INTO public.has_share_notes (id, user_id, to_user_id, seq) VALUES ('5524b99b82d7216d0e5516c0', '5368c1aa99c37b029d000001', '5524b99b99c37b2920000002', 0) ON CONFLICT DO NOTHING;
INSERT INTO public.has_share_notes (id, user_id, to_user_id, seq) VALUES ('5524ba2f82d7216d0e5516c4', '5368c1aa99c37b029d000001', '5524ba2f99c37b2920000007', 0) ON CONFLICT DO NOTHING;


--
-- Data for Name: note_content_histories; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.note_content_histories (id, user_id, histories) VALUES ('5368924119807a6ebd000000', '53688e1799c37b6889000001', '[{"Content": "<p>adfadfasdddd</p>", "UpdatedTime": "2014-05-06T16:14:24.855+08:00", "UpdatedUserId": "53688e1799c37b6889000001"}]') ON CONFLICT DO NOTHING;
INSERT INTO public.note_content_histories (id, user_id, histories) VALUES ('53689b9619807a1d82000000', '53688e1799c37b6889000001', '[{"Content": "you kown it\nadfadfad\nadfads\n\n\nadadfasdf", "UpdatedTime": "2014-05-06T16:22:13.541+08:00", "UpdatedUserId": "53688e1799c37b6889000001"}]') ON CONFLICT DO NOTHING;
INSERT INTO public.note_content_histories (id, user_id, histories) VALUES ('5368c1b919807a6f95000000', '5368c1aa99c37b029d000001', '[{"Content": "<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:&nbsp;</p><ul><li>管理知识</li><li>将知识公开成博客</li><li>分享知识</li><li>与好友协作知识<br></li><li>.........等待你来发现<br></li></ul><p>gemsnote已开源, 欢迎加入我们, 让gemsnote更好!</p><p>github:&nbsp;<a href=\"https://github.com/gemsnote/gemsnote\" data-mce-href=\"https://github.com/gemsnote/gemsnote\">https://github.com/gemsnote/gemsnote</a></p><p>gemsnote官网: <a href=\"http://gemsnote.com\" data-mce-href=\"http://gemsnote.com\">http://gemsnote.com</a></p><p>gemsnote官方博客:&nbsp;<a href=\"http://gemsnote.gemsnote.com\" data-mce-href=\"http://gemsnote.gemsnote.com\">http://gemsnote.gemsnote.com</a></p><p>lea++, gemsnote博客平台: <a href=\"http://lea.gemsnote.com\" data-mce-href=\"http://lea.gemsnote.com\">http://lea.gemsnote.com</a></p><p><img src=\"http://gemsnote.com/images/logo/gemsnote.png\" alt=\"\" width=\"230\" data-mce-src=\"http://gemsnote.com/images/logo/gemsnote.png\"></p>", "UpdatedTime": "2014-10-22T20:01:05.693+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:&nbsp;</p><ul><li>管理知识</li><li>将知识公开成博客</li><li>分享知识</li><li>与好友协作知识<br></li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>gemsnote已开源, 欢迎加入我们, 让gemsnote更好!</p><p>github:&nbsp;<a href=\"https://github.com/gemsnote/gemsnote\" data-mce-href=\"https://github.com/gemsnote/gemsnote\">https://github.com/gemsnote/gemsnote</a></p><p>官网: <a href=\"http://gemsnote.com\" data-mce-href=\"http://gemsnote.com\">http://gemsnote.com</a></p><p><img src=\"http://gemsnote.com/images/logo/gemsnote.png\" alt=\"\" width=\"230\" data-mce-src=\"http://gemsnote.com/images/logo/gemsnote.png\"></p>", "UpdatedTime": "2014-10-22T19:56:06.689+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:&nbsp;</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>gemsnote已开源, 欢迎加入我们, 让gemsnote更好!</p><p>github:&nbsp;<a href=\"https://github.com/gemsnote/gemsnote\" data-mce-href=\"https://github.com/gemsnote/gemsnote\">https://github.com/gemsnote/gemsnote</a></p><p>官网: <a href=\"http://gemsnote.com\" data-mce-href=\"http://gemsnote.com\">http://gemsnote.com</a></p><p><img src=\"http://gemsnote.com/images/logo/gemsnote.png\" alt=\"\" width=\"230\" data-mce-src=\"http://gemsnote.com/images/logo/gemsnote.png\"></p>", "UpdatedTime": "2014-09-04T15:41:22.443+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:&nbsp;</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>gemsnote已开源, 欢迎加入我们, 让gemsnote更好!</p><p>github:&nbsp;<a href=\"https://github.com/gemsnote/gemsnote\" data-mce-href=\"https://github.com/gemsnote/gemsnote\">https://github.com/gemsnote/gemsnote</a></p><p>官网: <a href=\"http://gemsnote.com\">http://gemsnote.com</a><a href=\"http://gemsnote.com\" data-mce-href=\"http://gemsnote.com\"></a></p><p><img src=\"http://gemsnote.com/images/logo/gemsnote.png\" alt=\"\" data-mce-src=\"http://gemsnote.com/images/logo/gemsnote.png\" width=\"230\"></p>", "UpdatedTime": "2014-09-04T15:41:10.837+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:&nbsp;</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>gemsnote已开源, 欢迎加入我们, 让gemsnote更好!</p><p>github:&nbsp;<a href=\"https://github.com/gemsnote/gemsnote\" data-mce-href=\"https://github.com/gemsnote/gemsnote\">https://github.com/gemsnote/gemsnote</a></p><p>官网: <a href=\"http://gemsnote.com\" data-mce-href=\"http://gemsnote.com\">http://gemsnote.com</a></p>", "UpdatedTime": "2014-05-06T19:10:35.445+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>gemsnote已开源, 欢迎加入我们, 让gemsnote更好!</p><p>github:&nbsp;<a href=\"https://github.com/gemsnote/gemsnote\" data-mce-href=\"https://github.com/gemsnote/gemsnote\">https://github.com/gemsnote/gemsnote</a></p><p>官网: <a href=\"http://gemsnote.com\" data-mce-href=\"http://gemsnote.com\">http://gemsnote.com</a></p>", "UpdatedTime": "2014-05-06T19:07:50.805+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>github:&nbsp;<a href=\"https://github.com/gemsnote/gemsnote\" data-mce-href=\"https://github.com/gemsnote/gemsnote\">https://github.com/gemsnote/gemsnote</a></p><p>官网: <a href=\"http://gemsnote.com\" data-mce-href=\"http://gemsnote.com\">http://gemsnote.com</a></p>", "UpdatedTime": "2014-05-06T19:06:56.064+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>github:&nbsp;<a href=\"https://github.com/gemsnote/gemsnote\">https://github.com/gemsnote/gemsnote</a></p><p>官网: <a href=\"http://gemsnote.com\">http://gemsnote.com</a></p>", "UpdatedTime": "2014-05-06T19:06:09.95+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>github:&nbsp;</p>", "UpdatedTime": "2014-05-06T19:05:49.321+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br data-mce-bogus=\"1\"></li></ul><p><br data-mce-bogus=\"1\"></p>", "UpdatedTime": "2014-05-06T19:05:39.672+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!</li></ul>", "UpdatedTime": "2014-05-06T19:05:27.511+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}]') ON CONFLICT DO NOTHING;
INSERT INTO public.note_content_histories (id, user_id, histories) VALUES ('5368c9fc99c37b095a00000a', '5368c9fc99c37b095a000006', '[{"Content": "<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:&nbsp;</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>gemsnote已开源, 欢迎加入我们, 让gemsnote更好!</p><p>github:&nbsp;<a href=\"https://github.com/gemsnote/gemsnote\" data-mce-href=\"https://github.com/gemsnote/gemsnote\">https://github.com/gemsnote/gemsnote</a></p><p>官网: <a href=\"http://gemsnote.com\" data-mce-href=\"http://gemsnote.com\">http://gemsnote.com</a></p><p><img src=\"http://gemsnote.com/images/logo/gemsnote.png\" alt=\"\" width=\"230\" data-mce-src=\"http://gemsnote.com/images/logo/gemsnote.png\"></p>", "UpdatedTime": "2014-09-04T15:40:01.282+08:00", "UpdatedUserId": "5368c9fc99c37b095a000006"}, {"Content": "<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:&nbsp;</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>gemsnote已开源, 欢迎加入我们, 让gemsnote更好!</p><p>github:&nbsp;<a href=\"https://github.com/gemsnote/gemsnote\" data-mce-href=\"https://github.com/gemsnote/gemsnote\">https://github.com/gemsnote/gemsnote</a></p><p>官网: <a href=\"http://gemsnote.com\">http://gemsnote.com</a><a href=\"http://gemsnote.com\" data-mce-href=\"http://gemsnote.com\"></a></p><p><img src=\"http://gemsnote.com/images/logo/gemsnote.png\" alt=\"\" data-mce-src=\"http://gemsnote.com/images/logo/gemsnote.png\" width=\"230\"></p>", "UpdatedTime": "2014-09-04T15:39:05.546+08:00", "UpdatedUserId": "5368c9fc99c37b095a000006"}, {"Content": "<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:&nbsp;</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>gemsnote已开源, 欢迎加入我们, 让gemsnote更好!</p><p>github:&nbsp;<a href=\"https://github.com/gemsnote/gemsnote\" data-mce-href=\"https://github.com/gemsnote/gemsnote\">https://github.com/gemsnote/gemsnote</a></p><p>官网: <a href=\"http://gemsnote.com\">http://gemsnote.com</a><a href=\"http://gemsnote.com\" data-mce-href=\"http://gemsnote.com\"></a></p><p><br data-mce-bogus=\"1\"></p>", "UpdatedTime": "2014-09-04T15:38:49.376+08:00", "UpdatedUserId": "5368c9fc99c37b095a000006"}, {"Content": "<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:&nbsp;</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>gemsnote已开源, 欢迎加入我们, 让gemsnote更好!</p><p>github:&nbsp;<a href=\"https://github.com/gemsnote/gemsnote\" data-mce-href=\"https://github.com/gemsnote/gemsnote\">https://github.com/gemsnote/gemsnote</a></p><p>官网: <a href=\"http://gemsnote.com\" data-mce-href=\"http://gemsnote.com\">http://gemsnote.com</a></p>", "UpdatedTime": "2014-09-04T15:38:46.891+08:00", "UpdatedUserId": "5368c9fc99c37b095a000006"}, {"Content": "<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:&nbsp;</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>gemsnote已开源, 欢迎加入我们, 让gemsnote更好!</p><p>github:&nbsp;<a href=\"https://github.com/gemsnote/gemsnote\" data-mce-href=\"https://github.com/gemsnote/gemsnote\">https://github.com/gemsnote/gemsnote</a></p><p>官网: <a href=\"http://gemsnote.com\">http://gemsnote.com</a></p><p><a href=\"http://gemsnote.com\" data-mce-href=\"http://gemsnote.com\"></a><img src=\"http://localhost/upload/5368c9fc99c37b095a000006/images/c93d600ee42c94876258cd0d51b36e06.png\" alt=\"\" data-mce-src=\"http://localhost/upload/5368c9fc99c37b095a000006/images/c93d600ee42c94876258cd0d51b36e06.png\" width=\"90\" style=\"line-height: 1.428571429;\"></p>", "UpdatedTime": "2014-09-04T15:31:04.11+08:00", "UpdatedUserId": "5368c9fc99c37b095a000006"}]') ON CONFLICT DO NOTHING;
INSERT INTO public.note_content_histories (id, user_id, histories) VALUES ('5447975919807a20ee000000', '5368c1aa99c37b029d000001', '[{"Content": "<p>dd</p>", "UpdatedTime": "2014-10-22T19:39:09.31+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}]') ON CONFLICT DO NOTHING;
INSERT INTO public.note_content_histories (id, user_id, histories) VALUES ('54632668a48d242f74000000', '5368c1aa99c37b029d000001', '[{"Content": "<p><br></p>", "UpdatedTime": "2014-11-12T17:20:58.51+08:00", "UpdatedUserId": "540817e099c37b583c000001"}]') ON CONFLICT DO NOTHING;
INSERT INTO public.note_content_histories (id, user_id, histories) VALUES ('5481218ff4e872105c000004', '5368c1aa99c37b029d000001', '[{"Content": "<p>java</p>", "UpdatedTime": "2014-12-05T11:29:40.654+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "<p><br></p>", "UpdatedTime": "2014-12-05T11:08:20.071+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}]') ON CONFLICT DO NOTHING;
INSERT INTO public.note_content_histories (id, user_id, histories) VALUES ('5481489cf4e8721ee3000000', '5368c1aa99c37b029d000001', '[{"Content": "<p><br></p>", "UpdatedTime": "2014-12-05T13:54:56.679+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}]') ON CONFLICT DO NOTHING;
INSERT INTO public.note_content_histories (id, user_id, histories) VALUES ('54817b33f4e8725881000000', '5368c1aa99c37b029d000001', '[{"Content": "<p><br></p>", "UpdatedTime": "2014-12-05T17:30:33.888+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}]') ON CONFLICT DO NOTHING;
INSERT INTO public.note_content_histories (id, user_id, histories) VALUES ('54817b49f4e8725881000001', '5368c1aa99c37b029d000001', '[{"Content": "<p><br></p>", "UpdatedTime": "2014-12-06T23:27:35.432+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "<p><br data-mce-bogus=\"1\"></p>", "UpdatedTime": "2014-12-06T23:27:32.358+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "<p>sefasfdasf</p>", "UpdatedTime": "2014-12-05T21:48:01.208+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "<p><br></p>", "UpdatedTime": "2014-12-05T17:30:53.255+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}]') ON CONFLICT DO NOTHING;
INSERT INTO public.note_content_histories (id, user_id, histories) VALUES ('5482b54cf4e87253cb000001', '5368c1aa99c37b029d000001', '[{"Content": "<p><br></p>", "UpdatedTime": "2014-12-06T15:50:43.7+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}]') ON CONFLICT DO NOTHING;
INSERT INTO public.note_content_histories (id, user_id, histories) VALUES ('54832064f4e87203a4000000', '5368c1aa99c37b029d000001', '[{"Content": "<p><br></p>", "UpdatedTime": "2014-12-06T23:27:56.648+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}]') ON CONFLICT DO NOTHING;
INSERT INTO public.note_content_histories (id, user_id, histories) VALUES ('5483207cf4e87203a4000001', '5368c1aa99c37b029d000001', '[{"Content": "\n# About Gemsnote\n\n## 1. Introduction\n\nGemsnote, not just a notepad!\n![gemsnote](http://7xj51o.com1.z0.glb.clouddn.com/default_markdown.png)\n\n**Some Features**\n\n* Knowledge: Manage your knowledge in gemsnote. gemsnote contains the tinymce editor and a markdown editor, just enjoy yourself writing.\n* Share: Share your knowledge with your friends in gemsnote. You can invite your friends to join your notepad in the cloud so you can share knowledge.\n* Cooperation: Collaborate with friends to improve your skills.\n* Blog: Publish your knowledge and make gemsnote your blog.\n\n## 2. Why we created gemsnote\nTo be honest, our inspiration comes from Evernote. We use Evernote to manage our knowledge everyday. But we find that:\n* Evernote''s editor can''t meet our needs, it does not have document navigation, it does not render code properly (as a programmer, syntax highlighted code rendering is a basic need), it cannot resize images and so forth\n* We like markdown, but Evernote does not support it.\n* We want to share our knowledge, so all of us have our blogs (e.g. on Wordpress) and our Evernote accounts, but why can not those two be one!\n* ......\n\n## 3. How to install gemsnote\n\nMore information about how to install gemsnote please see:\n* [gemsnote binary distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-binary-distribution-installation-tutorial)\n* [gemsnote develop distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-develop-distribution-installation-tutorial)\n\n## 4. How to develop gemsnote\n\nPlease see [How-to-develop-gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91gemsnote)\n\n\n## 5. Docs\n* [gemsnote binary distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-binary-distribution-installation-tutorial)\n* [gemsnote develop distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-develop-distribution-installation-tutorial)\n* [gemsnote blog theme api](https://github.com/gemsnote/gemsnote/wiki/gemsnote-blog-theme-api_en)\n\nMore docs please see [wiki](https://github.com/gemsnote/gemsnote/wiki).\n\n## 6. Contributors\nThank you to all the [contributors](https://github.com/gemsnote/gemsnote/graphs/contributors) on\nthis project. Your help is much appreciated.\n\n## 7.Join us\nPlease fork this repository and contribute back using [pull requests](https://github.com/gemsnote/gemsnote/pulls).\n\nIf you find some problems or has some good ideas, please submit [issues](https://github.com/gemsnote/gemsnote/issues).\n\nYou are always welcomed!\n\n## 8. Donation\nSupport us, [donate us](http://gemsnote.org/#donate). And thanks [donators](http://gemsnote.gemsnote.com/post/gemsnote-donation-list).\n\n## 9. Related projects\n* [Gemsnote Desktop App](https://github.com/gemsnote/desktop-app), [Download](http://app.gemsnote.com)\n* [Gemsnote IOS](https://github.com/gemsnote/gemsnote-ios), development phase\n* [Gemsnote Android](https://github.com/Dminter/gemsnote-android-client), development phase\n\nAnd also, you are welcome to join us.\n\n## 9. Discussion\n* [gemsnote bbs](http://bbs.gemsnote.com)\n* [gemsnote google group](https://groups.google.com/forum/#!forum/gemsnote)\n* QQ Group: 158716820\n\n\n", "UpdatedTime": "2015-06-15T18:37:30.724+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "\n# About Gemsnote\n\n## 1. Introduction\n\nGemsnote, not just a notepad!\n![gemsnote](http://7xj51o.com1.z0.glb.clouddn.com/default_markdown.png)\n\n**Some Features**\n\n* Knowledge: Manage your knowledge in gemsnote. gemsnote contains the tinymce editor and a markdown editor, just enjoy yourself writing.\n* Share: Share your knowledge with your friends in gemsnote. You can invite your friends to join your notepad in the cloud so you can share knowledge.\n* Cooperation: Collaborate with friends to improve your skills.\n* Blog: Publish your knowledge and make gemsnote your blog.\n\n## 2. Why we created gemsnote\nTo be honest, our inspiration comes from Evernote. We use Evernote to manage our knowledge everyday. But we find that:\n* Evernote''s editor can''t meet our needs, it does not have document navigation, it does not render code properly (as a programmer, syntax highlighted code rendering is a basic need), it cannot resize images and so forth\n* We like markdown, but Evernote does not support it.\n* We want to share our knowledge, so all of us have our blogs (e.g. on Wordpress) and our Evernote accounts, but why can not those two be one!\n* ......\n\n## 3. How to install gemsnote\n\nMore information about how to install gemsnote please see:\n* [gemsnote binary distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-binary-distribution-installation-tutorial)\n* [gemsnote develop distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-develop-distribution-installation-tutorial)\n\n## 4. How to develop gemsnote\n\nPlease see [How-to-develop-gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91gemsnote)\n\n\n## 5. Docs\n* [gemsnote binary distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-binary-distribution-installation-tutorial)\n* [gemsnote develop distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-develop-distribution-installation-tutorial)\n* [gemsnote blog theme api](https://github.com/gemsnote/gemsnote/wiki/gemsnote-blog-theme-api_en)\n\nMore docs please see [wiki](https://github.com/gemsnote/gemsnote/wiki).\n\n## 6. Contributors\nThank you to all the [contributors](https://github.com/gemsnote/gemsnote/graphs/contributors) on\nthis project. Your help is much appreciated.\n\n## 7.Join us\nPlease fork this repository and contribute back using [pull requests](https://github.com/gemsnote/gemsnote/pulls).\n\nIf you find some problems or has some good ideas, please submit [issues](https://github.com/gemsnote/gemsnote/issues).\n\nYou are always welcomed!\n\n## 8. Donation\nSupport us, [donate us](http://gemsnote.org/#donate). And thanks [donators](http://gemsnote.gemsnote.com/post/gemsnote-donation-list).\n\n## 9. Related projects\n* [Gemsnote Desktop App](https://github.com/gemsnote/desktop-app), [Download](http://app.gemsnote.com)\n* [Gemsnote IOS](https://github.com/gemsnote/gemsnote-ios), development phase\n* [Gemsnote Android](https://github.com/Dminter/gemsnote-android-client), development phase\n* \n\nAnd also, you are welcome to join us.\n\n## 9. Discussion\n* [gemsnote bbs](http://bbs.gemsnote.com)\n* [gemsnote google group](https://groups.google.com/forum/#!forum/gemsnote)\n* QQ Group: 158716820\n\n\n", "UpdatedTime": "2015-06-15T18:37:29.098+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "\n# About Gemsnote\n\n## 1. Introduction\n\nGemsnote, not just a notepad!\n![gemsnote](http://7xj51o.com1.z0.glb.clouddn.com/default_markdown.png)\n\n**Some Features**\n\n* Knowledge: Manage your knowledge in gemsnote. gemsnote contains the tinymce editor and a markdown editor, just enjoy yourself writing.\n* Share: Share your knowledge with your friends in gemsnote. You can invite your friends to join your notepad in the cloud so you can share knowledge.\n* Cooperation: Collaborate with friends to improve your skills.\n* Blog: Publish your knowledge and make gemsnote your blog.\n\n## 2. Why we created gemsnote\nTo be honest, our inspiration comes from Evernote. We use Evernote to manage our knowledge everyday. But we find that:\n* Evernote''s editor can''t meet our needs, it does not have document navigation, it does not render code properly (as a programmer, syntax highlighted code rendering is a basic need), it cannot resize images and so forth\n* We like markdown, but Evernote does not support it.\n* We want to share our knowledge, so all of us have our blogs (e.g. on Wordpress) and our Evernote accounts, but why can not those two be one!\n* ......\n\n## 3. How to install gemsnote\n\nMore information about how to install gemsnote please see:\n* [gemsnote binary distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-binary-distribution-installation-tutorial)\n* [gemsnote develop distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-develop-distribution-installation-tutorial)\n\n## 4. How to develop gemsnote\n\nPlease see [How-to-develop-gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91gemsnote)\n\n\n## 5. Docs\n* [gemsnote binary distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-binary-distribution-installation-tutorial)\n* [gemsnote develop distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-develop-distribution-installation-tutorial)\n* [gemsnote blog theme api](https://github.com/gemsnote/gemsnote/wiki/gemsnote-blog-theme-api_en)\n\nMore docs please see [wiki](https://github.com/gemsnote/gemsnote/wiki).\n\n## 6. Contributors\nThank you to all the [contributors](https://github.com/gemsnote/gemsnote/graphs/contributors) on\nthis project. Your help is much appreciated.\n\n## 7.Join us\nPlease fork this repository and contribute back using [pull requests](https://github.com/gemsnote/gemsnote/pulls).\n\nIf you find some problems or has some good ideas, please submit [issues](https://github.com/gemsnote/gemsnote/issues).\n\nYou are always welcomed!\n\n## 8. Donation\nSupport us, [donate us](http://gemsnote.org/#donate). And thanks [donators](http://gemsnote.gemsnote.com/post/gemsnote-donation-list).\n\n## 9. Related projects\n* [Gemsnote Desktop App](https://github.com/gemsnote/desktop-app), [Download](http://app.gemsnote.com)\n* [Gemsnote IOS](https://github.com/gemsnote/gemsnote-ios), development phase\n* [Gemsnote Android](https://github.com/Dminter/gemsnote-android-client), development phase\n\nAnd also, you are welcome to join us.\n\n## 9. Discussion\n* [gemsnote bbs](http://bbs.gemsnote.com)\n* [gemsnote google group](https://groups.google.com/forum/#!forum/gemsnote)\n* QQ Group: 158716820\n\n\n", "UpdatedTime": "2015-06-15T18:37:16.05+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "\n# About Gemsnote\n\n## 1. Introduction\n\nGemsnote, not just a notepad!\n![gemsnote](http://7xj51o.com1.z0.glb.clouddn.com/default_markdown.png)\n\n**Some Features**\n\n* Knowledge: Manage your knowledge in gemsnote. gemsnote contains the tinymce editor and a markdown editor, just enjoy yourself writing.\n* Share: Share your knowledge with your friends in gemsnote. You can invite your friends to join your notepad in the cloud so you can share knowledge.\n* Cooperation: Collaborate with friends to improve your skills.\n* Blog: Publish your knowledge and make gemsnote your blog.\n\n## 2. Why we created gemsnote\nTo be honest, our inspiration comes from Evernote. We use Evernote to manage our knowledge everyday. But we find that:\n* Evernote''s editor can''t meet our needs, it does not have document navigation, it does not render code properly (as a programmer, syntax highlighted code rendering is a basic need), it cannot resize images and so forth\n* We like markdown, but Evernote does not support it.\n* We want to share our knowledge, so all of us have our blogs (e.g. on Wordpress) and our Evernote accounts, but why can not those two be one!\n* ......\n\n## 3. How to install gemsnote\n\nMore information about how to install gemsnote please see:\n* [gemsnote binary distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-binary-distribution-installation-tutorial)\n* [gemsnote develop distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-develop-distribution-installation-tutorial)\n\n## 4. How to develop gemsnote\n\nPlease see [How-to-develop-gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91gemsnote)\n\n\n## 5. Docs\n* [gemsnote binary distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-binary-distribution-installation-tutorial)\n* [gemsnote develop distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-develop-distribution-installation-tutorial)\n* [gemsnote blog theme api](https://github.com/gemsnote/gemsnote/wiki/gemsnote-blog-theme-api_en)\n\nMore docs please see [wiki](https://github.com/gemsnote/gemsnote/wiki).\n\n## 6. Contributors\nThank you to all the [contributors](https://github.com/gemsnote/gemsnote/graphs/contributors) on\nthis project. Your help is much appreciated.\n\n## 7.Join us\nPlease fork this repository and contribute back using [pull requests](https://github.com/gemsnote/gemsnote/pulls).\n\nIf you find some problems or has some good ideas, please submit [issues](https://github.com/gemsnote/gemsnote/issues).\n\nYou are always welcomed!\n\n## 8. Donation\nSupport us, [donate us](http://gemsnote.org/#donate). And thanks [donators](http://gemsnote.gemsnote.com/post/gemsnote-donation-list).\n\n## 9. Related projects\n* [Gemsnote Desktop App](https://github.com/gemsnote/desktop-app), [Download](http://app.gemsnote.com)\n* [Gemsnote IOS](https://github.com/gemsnote/gemsnote-ios), development phase\n* [Gemsnote Android](https://github.com/Dminter/gemsnote-android-client), development phase\n\nAnd also, you are welcome to join us.\n\n## 9. Discussion\n* [gemsnote bbs](http://bbs.gemsnote.com)\n* [gemsnote google group](https://groups.google.com/forum/#!forum/gemsnote)\n* QQ Group: 158716820\n\n-----------------------------------------------------------------------\n[中文](README_zh.md)\n\n\n", "UpdatedTime": "2015-06-15T18:36:59.145+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "\n# Gemsnote\n\n## 1. Introduction\n\nGemsnote, not just a notepad!\n![gemsnote](http://7xj51o.com1.z0.glb.clouddn.com/default_markdown.png)\n\n**Some Features**\n\n* Knowledge: Manage your knowledge in gemsnote. gemsnote contains the tinymce editor and a markdown editor, just enjoy yourself writing.\n* Share: Share your knowledge with your friends in gemsnote. You can invite your friends to join your notepad in the cloud so you can share knowledge.\n* Cooperation: Collaborate with friends to improve your skills.\n* Blog: Publish your knowledge and make gemsnote your blog.\n\n## 2. Why we created gemsnote\nTo be honest, our inspiration comes from Evernote. We use Evernote to manage our knowledge everyday. But we find that:\n* Evernote''s editor can''t meet our needs, it does not have document navigation, it does not render code properly (as a programmer, syntax highlighted code rendering is a basic need), it cannot resize images and so forth\n* We like markdown, but Evernote does not support it.\n* We want to share our knowledge, so all of us have our blogs (e.g. on Wordpress) and our Evernote accounts, but why can not those two be one!\n* ......\n\n## 3. How to install gemsnote\n\nMore information about how to install gemsnote please see:\n* [gemsnote binary distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-binary-distribution-installation-tutorial)\n* [gemsnote develop distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-develop-distribution-installation-tutorial)\n\n## 4. How to develop gemsnote\n\nPlease see [How-to-develop-gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91gemsnote)\n\n\n## 5. Docs\n* [gemsnote binary distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-binary-distribution-installation-tutorial)\n* [gemsnote develop distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-develop-distribution-installation-tutorial)\n* [gemsnote blog theme api](https://github.com/gemsnote/gemsnote/wiki/gemsnote-blog-theme-api_en)\n\nMore docs please see [wiki](https://github.com/gemsnote/gemsnote/wiki).\n\n## 6. Contributors\nThank you to all the [contributors](https://github.com/gemsnote/gemsnote/graphs/contributors) on\nthis project. Your help is much appreciated.\n\n## 7.Join us\nPlease fork this repository and contribute back using [pull requests](https://github.com/gemsnote/gemsnote/pulls).\n\nIf you find some problems or has some good ideas, please submit [issues](https://github.com/gemsnote/gemsnote/issues).\n\nYou are always welcomed!\n\n## 8. Donation\nSupport us, [donate us](http://gemsnote.org/#donate). And thanks [donators](http://gemsnote.gemsnote.com/post/gemsnote-donation-list).\n\n## 9. Related projects\n* [Gemsnote Desktop App](https://github.com/gemsnote/desktop-app), [Download](http://app.gemsnote.com)\n* [Gemsnote IOS](https://github.com/gemsnote/gemsnote-ios), development phase\n* [Gemsnote Android](https://github.com/Dminter/gemsnote-android-client), development phase\n\nAnd also, you are welcome to join us.\n\n## 9. Discussion\n* [gemsnote bbs](http://bbs.gemsnote.com)\n* [gemsnote google group](https://groups.google.com/forum/#!forum/gemsnote)\n* QQ Group: 158716820\n\n-----------------------------------------------------------------------\n[中文](README_zh.md)\n\n\n", "UpdatedTime": "2015-06-15T18:35:37.288+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "\n# Gemsnote\n\n## 1. Introduction\n\nGemsnote, not just a notepad!\n![gemsnote.png](gemsnote.png \"\")\n\n**Some Features**\n\n* Knowledge: Manage your knowledge in gemsnote. gemsnote contains the tinymce editor and a markdown editor, just enjoy yourself writing.\n* Share: Share your knowledge with your friends in gemsnote. You can invite your friends to join your notepad in the cloud so you can share knowledge.\n* Cooperation: Collaborate with friends to improve your skills.\n* Blog: Publish your knowledge and make gemsnote your blog.\n\n## 2. Why we created gemsnote\nTo be honest, our inspiration comes from Evernote. We use Evernote to manage our knowledge everyday. But we find that:\n* Evernote''s editor can''t meet our needs, it does not have document navigation, it does not render code properly (as a programmer, syntax highlighted code rendering is a basic need), it cannot resize images and so forth\n* We like markdown, but Evernote does not support it.\n* We want to share our knowledge, so all of us have our blogs (e.g. on Wordpress) and our Evernote accounts, but why can not those two be one!\n* ......\n\n## 3. How to install gemsnote\n\nMore information about how to install gemsnote please see:\n* [gemsnote binary distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-binary-distribution-installation-tutorial)\n* [gemsnote develop distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-develop-distribution-installation-tutorial)\n\n## 4. How to develop gemsnote\n\nPlease see [How-to-develop-gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91gemsnote)\n\n\n## 5. Docs\n* [gemsnote binary distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-binary-distribution-installation-tutorial)\n* [gemsnote develop distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-develop-distribution-installation-tutorial)\n* [gemsnote blog theme api](https://github.com/gemsnote/gemsnote/wiki/gemsnote-blog-theme-api_en)\n\nMore docs please see [wiki](https://github.com/gemsnote/gemsnote/wiki).\n\n## 6. Contributors\nThank you to all the [contributors](https://github.com/gemsnote/gemsnote/graphs/contributors) on\nthis project. Your help is much appreciated.\n\n## 7.Join us\nPlease fork this repository and contribute back using [pull requests](https://github.com/gemsnote/gemsnote/pulls).\n\nIf you find some problems or has some good ideas, please submit [issues](https://github.com/gemsnote/gemsnote/issues).\n\nYou are always welcomed!\n\n## 8. Donation\nSupport us, [donate us](http://gemsnote.org/#donate). And thanks [donators](http://gemsnote.gemsnote.com/post/gemsnote-donation-list).\n\n## 9. Related projects\n* [Gemsnote Desktop App](https://github.com/gemsnote/desktop-app), [Download](http://app.gemsnote.com)\n* [Gemsnote IOS](https://github.com/gemsnote/gemsnote-ios), development phase\n* [Gemsnote Android](https://github.com/Dminter/gemsnote-android-client), development phase\n\nAnd also, you are welcome to join us.\n\n## 9. Discussion\n* [gemsnote bbs](http://bbs.gemsnote.com)\n* [gemsnote google group](https://groups.google.com/forum/#!forum/gemsnote)\n* QQ Group: 158716820\n\n-----------------------------------------------------------------------\n[中文](README_zh.md)\n\n\n", "UpdatedTime": "2015-06-15T18:32:38.663+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "# Gemsnote产品说明\n\n## 1. 介绍\n\nGemsnote, 不只是笔记!\n\n**特性**\n\n* 知识管理: 通过gemsnote来管理知识, gemsnote有易操作的界面, 包含两款编辑器tinymce和markdown. 在gemsnote, 你可以尽情享受写作.\n* 分享: 你也可以通过分享知识给好友, 让好友拥有你的知识.\n* 协作: 在分享的同时也可以与好友一起协作知识.\n* 博客: gemsnote也可以作为你的博客, 将知识公开成博客, 让gemsnote把你的知识传播的更远!\n\n## 2. 为什么我们要创建gemsnote?\n说实话, 我们曾是evernote的忠实粉丝, 但是我们也发现evernote的不足:\n* evernote的编辑器不能满足我们的需求, 不能贴代码(格式会乱掉, 作为程序员, 代码是我们的基本需求啊), 图片不能缩放.\n* 我们是markdown的爱好者, 可是evernote竟然没有.\n* 我们也想将知识公开, 所以我们有自己的博客, 如wordpress, 但为什么这两者不能合二为一呢?\n* 还有...\n\n## 3.安装gemsnote\ngemsnote是一款私有云笔记, 你可以下载它安装在自己的服务器上, 当然也可以在 http://gemsnote.com 上注册.\n\n这里详细整理了gemsnote二进版和gemsnote开发版的安装教程, 请移步至:\n* [gemsnote二进制详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n* [gemsnote开发版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n\n### 3.1. 下载gemsnote\n\nGemsnote V1.0-beta 已发布, 二进制文件(暂时没有windows版的):\n\n* Linux: [gemsnote-linux-x86_64.v1.0-beta.bin.tar.gz](https://github.com/gemsnote/gemsnote/releases/download/1.0-beta/gemsnote-linux-x86_64.v1.0-beta.bin.tar.gz)\n* MacOS X: [gemsnote-mac-x86_64.v1.0-beta.bin.tar.gz](https://github.com/gemsnote/gemsnote/releases/download/1.0-beta/gemsnote-mac-x86_64.v1.0-beta.bin.tar.gz)\n\n### 3.2. 安装 MongodbDB\n\nGemsnote是由golang(使用[revel](https://revel.github.io/)框架 和 [MongoDB](https://www.mongodb.org)数据库), 你需要先安装Mongodb.\n\n安装MongodbDB, 导入数据更多细节请查看: [wiki](https://github.com/gemsnote/gemsnote/wiki/Install-Mongodb)\n\n### 3.3. 导入初始数据\n\nMongodbDB初始数据在 `[PATH_TO_GEMSNOTE]/mongodb_backup/gemsnote_install_data`\n\n```\n$> mongorestore -h localhost -d gemsnote --directoryperdb PATH_TO_GEMSNOTE/mongodb_backup/gemsnote_install_data\n```\n\n初始数据包含两个用户:\n\n```\nuser2 username: admin, password: abc123 (管理员, 重要!)\nuser3 username: demo@gemsnote.com, password: demo@gemsnote.com (为体验使用)\n```\n\n### 3.4. 配置\n\n修改 `[PATH_TO_GEMSNOTE]/conf/app.conf`. 有以下选项:\n\n``mongodb``  **必须配置!**\n\n```Shell\ndb.host=localhost\ndb.port=27017\ndb.dbname=gemsnote\ndb.username=\ndb.password=\n```\n\n``app.secret`` **重要**\n请随意修改一个, app的密钥, 不能使用默认的, 不然会有安全问题\n\n更多配置请查看 `app/app.conf` 和 [revel 手册](https://revel.github.io/)\n\n### 3.5. 运行gemsnote\n\n```\n$> cd PATH_TO_GEMSNOTE/bin\n$> sudo sh run.sh\n```\n\n## 4. 如何对gemsnote进行二次开发\n\n请查看 [How-to-develop-gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91gemsnote)\n\n## 5. 贡献者\n多谢 [贡献者](https://github.com/gemsnote/gemsnote/graphs/contributors) 的贡献, gemsnote因有你们而更完美!\n\n## 6. 加入我们\n\n欢迎提交[pull requests](https://github.com/gemsnote/gemsnote/pulls) 到gemsnote.\n\ngemsnote还有很多问题, 如果你喜欢它, 欢迎加入我们一起完善gemsnote.\n\n## 讨论\n* [gemsnote 社区](http://bbs.gemsnote.com)\n* QQ群: 158716820\n* [gemsnote google group](https://groups.google.com/forum/#!forum/gemsnote)", "UpdatedTime": "2014-12-06T23:28:40.476+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}]') ON CONFLICT DO NOTHING;
INSERT INTO public.note_content_histories (id, user_id, histories) VALUES ('557eaa9905fcd14d95000001', '5368c1aa99c37b029d000001', '[{"Content": "# 关于Gemsnote\n\n## 1. 介绍\n\nGemsnote, 不只是笔记!\n![gemsnote](http://7xj51o.com1.z0.glb.clouddn.com/default_markdown.png)\n\n**特性**\n\n* 知识管理: 通过gemsnote来管理知识, gemsnote有易操作的界面, 包含两款编辑器tinymce和markdown. 在gemsnote, 你可以尽情享受写作.\n* 分享: 你也可以通过分享知识给好友, 让好友拥有你的知识.\n* 协作: 在分享的同时也可以与好友一起协作知识.\n* 博客: gemsnote也可以作为你的博客, 将知识公开成博客, 让gemsnote把你的知识传播的更远!\n\n## 2. 为什么我们要创建gemsnote?\n说实话, 我们曾是evernote的忠实粉丝, 但是我们也发现evernote的不足:\n* evernote的编辑器不能满足我们的需求, 不能贴代码(格式会乱掉, 作为程序员, 代码是我们的基本需求啊), 图片不能缩放.\n* 我们是markdown的爱好者, 可是evernote竟然没有.\n* 我们也想将知识公开, 所以我们有自己的博客, 如wordpress, 但为什么这两者不能合二为一呢?\n* 还有...\n\n## 3.安装gemsnote\ngemsnote是一款私有云笔记, 你可以下载它安装在自己的服务器上, 当然也可以在 http://gemsnote.com 上注册.\n\n这里详细整理了gemsnote二进版和gemsnote开发版的安装教程, 请移步至:\n\n* [gemsnote二进制详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n* [gemsnote开发版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n\n## 4. 如何对gemsnote进行二次开发\n\n请查看 [How-to-develop-gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91gemsnote)\n\n## 5 相关文档\n* [gemsnote二进制版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n* [gemsnote开发版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n* [Gemsnote source gemsnote源码导读](https://github.com/gemsnote/gemsnote/wiki/Gemsnote-source-gemsnote源码导读)\n* [gemsnote blog theme api(中文版)](https://github.com/gemsnote/gemsnote/wiki/gemsnote-blog-theme-api)\n* [How to develop gemsnote 如何开发gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-如何开发gemsnote)\n\n更多文档请查看 [wiki](https://github.com/gemsnote/gemsnote/wiki).\n\n## 6. 贡献者\n多谢 [贡献者](https://github.com/gemsnote/gemsnote/graphs/contributors) 的贡献, gemsnote因有你们而更完美!\n\n## 7. 加入我们\n\n欢迎提交[pull requests](https://github.com/gemsnote/gemsnote/pulls) 到gemsnote.\n\n有任何问题或建议, 欢迎提交[issue](https://github.com/gemsnote/gemsnote/issues).\n\nGemsnote还有很多问题, 如果你喜欢它, 欢迎加入我们一起完善gemsnote.\n\n## 8. 捐赠\n支持我们, [捐赠Gemsnote](http://gemsnote.org/#donate). 感谢[捐赠者](http://gemsnote.gemsnote.com/post/gemsnote-donation-list), 谢谢你们的鼓励, Gemsnote会一直坚持!\n\n## 9. 其它相关项目\n* [Gemsnote Desktop App](https://github.com/gemsnote/desktop-app), [下载地址](http://app.gemsnote.com)\n* [Gemsnote IOS](https://github.com/gemsnote/gemsnote-ios), 开发阶段\n* [Gemsnote Android](https://github.com/Dminter/gemsnote-android-client), 开发阶段\n\n同样, 欢迎加入我们!\n\n## 讨论\n* [gemsnote 社区](http://bbs.gemsnote.com)\n* QQ群: 158716820\n* [gemsnote google group](https://groups.google.com/forum/#!forum/gemsnote)", "UpdatedTime": "2015-06-15T18:37:24.15+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "# 关于Gemsnote\n\n## 1. 介绍\n\nGemsnote, 不只是笔记!\n![gemsnote](http://7xj51o.com1.z0.glb.clouddn.com/default_markdown.png)\n\n**特性**\n\n* 知识管理: 通过gemsnote来管理知识, gemsnote有易操作的界面, 包含两款编辑器tinymce和markdown. 在gemsnote, 你可以尽情享受写作.\n* 分享: 你也可以通过分享知识给好友, 让好友拥有你的知识.\n* 协作: 在分享的同时也可以与好友一起协作知识.\n* 博客: gemsnote也可以作为你的博客, 将知识公开成博客, 让gemsnote把你的知识传播的更远!\n\n## 2. 为什么我们要创建gemsnote?\n说实话, 我们曾是evernote的忠实粉丝, 但是我们也发现evernote的不足:\n* evernote的编辑器不能满足我们的需求, 不能贴代码(格式会乱掉, 作为程序员, 代码是我们的基本需求啊), 图片不能缩放.\n* 我们是markdown的爱好者, 可是evernote竟然没有.\n* 我们也想将知识公开, 所以我们有自己的博客, 如wordpress, 但为什么这两者不能合二为一呢?\n* 还有...\n\n## 3.安装gemsnote\ngemsnote是一款私有云笔记, 你可以下载它安装在自己的服务器上, 当然也可以在 http://gemsnote.com 上注册.\n\n这里详细整理了gemsnote二进版和gemsnote开发版的安装教程, 请移步至:\n\n* [gemsnote二进制详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n* [gemsnote开发版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n\n## 4. 如何对gemsnote进行二次开发\n\n请查看 [How-to-develop-gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91gemsnote)\n\n## 5 相关文档\n* [gemsnote二进制版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n* [gemsnote开发版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n* [Gemsnote source gemsnote源码导读](https://github.com/gemsnote/gemsnote/wiki/Gemsnote-source-gemsnote源码导读)\n* [gemsnote blog theme api(中文版)](https://github.com/gemsnote/gemsnote/wiki/gemsnote-blog-theme-api)\n* [How to develop gemsnote 如何开发gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-如何开发gemsnote)\n\n更多文档请查看 [wiki](https://github.com/gemsnote/gemsnote/wiki).\n\n## 6. 贡献者\n多谢 [贡献者](https://github.com/gemsnote/gemsnote/graphs/contributors) 的贡献, gemsnote因有你们而更完美!\n\n## 7. 加入我们\n\n欢迎提交[pull requests](https://github.com/gemsnote/gemsnote/pulls) 到gemsnote.\n\n有任何问题或建议, 欢迎提交[issue](https://github.com/gemsnote/gemsnote/issues).\n\nGemsnote还有很多问题, 如果你喜欢它, 欢迎加入我们一起完善gemsnote.\n\n## 8. 捐赠\n支持我们, [捐赠Gemsnote](http://gemsnote.org/#donate). 感谢[捐赠者](http://gemsnote.gemsnote.com/post/gemsnote-donation-list), 谢谢你们的鼓励, Gemsnote会一直坚持!\n\n## 9. 其它相关项目\n* [Gemsnote Desktop App](https://github.com/gemsnote/desktop-app), [下载地址](http://app.gemsnote.com)\n* [Gemsnote IOS](https://github.com/gemsnote/gemsnote-ios), 开发阶段\n* [Gemsnote Android](https://github.com/Dminter/gemsnote-android-client), 开发阶段\n\n同样, 欢迎加入我们!\n\n## 讨论\n* [gemsnote 社区](http://bbs.gemsnote.com)\n* QQ群: 158716820\n* [gemsnote google group](https://groups.google.com/forum/#!forum/gemsnote)\n\n----------------------------------------------------------------\n[English](README.md)", "UpdatedTime": "2015-06-15T18:36:51.151+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "# Gemsnote\n\n## 1. 介绍\n\nGemsnote, 不只是笔记!\n![gemsnote](http://7xj51o.com1.z0.glb.clouddn.com/default_markdown.png)\n\n**特性**\n\n* 知识管理: 通过gemsnote来管理知识, gemsnote有易操作的界面, 包含两款编辑器tinymce和markdown. 在gemsnote, 你可以尽情享受写作.\n* 分享: 你也可以通过分享知识给好友, 让好友拥有你的知识.\n* 协作: 在分享的同时也可以与好友一起协作知识.\n* 博客: gemsnote也可以作为你的博客, 将知识公开成博客, 让gemsnote把你的知识传播的更远!\n\n## 2. 为什么我们要创建gemsnote?\n说实话, 我们曾是evernote的忠实粉丝, 但是我们也发现evernote的不足:\n* evernote的编辑器不能满足我们的需求, 不能贴代码(格式会乱掉, 作为程序员, 代码是我们的基本需求啊), 图片不能缩放.\n* 我们是markdown的爱好者, 可是evernote竟然没有.\n* 我们也想将知识公开, 所以我们有自己的博客, 如wordpress, 但为什么这两者不能合二为一呢?\n* 还有...\n\n## 3.安装gemsnote\ngemsnote是一款私有云笔记, 你可以下载它安装在自己的服务器上, 当然也可以在 http://gemsnote.com 上注册.\n\n这里详细整理了gemsnote二进版和gemsnote开发版的安装教程, 请移步至:\n\n* [gemsnote二进制详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n* [gemsnote开发版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n\n## 4. 如何对gemsnote进行二次开发\n\n请查看 [How-to-develop-gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91gemsnote)\n\n## 5 相关文档\n* [gemsnote二进制版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n* [gemsnote开发版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n* [Gemsnote source gemsnote源码导读](https://github.com/gemsnote/gemsnote/wiki/Gemsnote-source-gemsnote源码导读)\n* [gemsnote blog theme api(中文版)](https://github.com/gemsnote/gemsnote/wiki/gemsnote-blog-theme-api)\n* [How to develop gemsnote 如何开发gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-如何开发gemsnote)\n\n更多文档请查看 [wiki](https://github.com/gemsnote/gemsnote/wiki).\n\n## 6. 贡献者\n多谢 [贡献者](https://github.com/gemsnote/gemsnote/graphs/contributors) 的贡献, gemsnote因有你们而更完美!\n\n## 7. 加入我们\n\n欢迎提交[pull requests](https://github.com/gemsnote/gemsnote/pulls) 到gemsnote.\n\n有任何问题或建议, 欢迎提交[issue](https://github.com/gemsnote/gemsnote/issues).\n\nGemsnote还有很多问题, 如果你喜欢它, 欢迎加入我们一起完善gemsnote.\n\n## 8. 捐赠\n支持我们, [捐赠Gemsnote](http://gemsnote.org/#donate). 感谢[捐赠者](http://gemsnote.gemsnote.com/post/gemsnote-donation-list), 谢谢你们的鼓励, Gemsnote会一直坚持!\n\n## 9. 其它相关项目\n* [Gemsnote Desktop App](https://github.com/gemsnote/desktop-app), [下载地址](http://app.gemsnote.com)\n* [Gemsnote IOS](https://github.com/gemsnote/gemsnote-ios), 开发阶段\n* [Gemsnote Android](https://github.com/Dminter/gemsnote-android-client), 开发阶段\n\n同样, 欢迎加入我们!\n\n## 讨论\n* [gemsnote 社区](http://bbs.gemsnote.com)\n* QQ群: 158716820\n* [gemsnote google group](https://groups.google.com/forum/#!forum/gemsnote)\n\n----------------------------------------------------------------\n[English](README.md)", "UpdatedTime": "2015-06-15T18:36:44.244+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "# Gemsnote产品\n\n## 1. 介绍\n\nGemsnote, 不只是笔记!\n![gemsnote](http://7xj51o.com1.z0.glb.clouddn.com/default_markdown.png)\n\n**特性**\n\n* 知识管理: 通过gemsnote来管理知识, gemsnote有易操作的界面, 包含两款编辑器tinymce和markdown. 在gemsnote, 你可以尽情享受写作.\n* 分享: 你也可以通过分享知识给好友, 让好友拥有你的知识.\n* 协作: 在分享的同时也可以与好友一起协作知识.\n* 博客: gemsnote也可以作为你的博客, 将知识公开成博客, 让gemsnote把你的知识传播的更远!\n\n## 2. 为什么我们要创建gemsnote?\n说实话, 我们曾是evernote的忠实粉丝, 但是我们也发现evernote的不足:\n* evernote的编辑器不能满足我们的需求, 不能贴代码(格式会乱掉, 作为程序员, 代码是我们的基本需求啊), 图片不能缩放.\n* 我们是markdown的爱好者, 可是evernote竟然没有.\n* 我们也想将知识公开, 所以我们有自己的博客, 如wordpress, 但为什么这两者不能合二为一呢?\n* 还有...\n\n## 3.安装gemsnote\ngemsnote是一款私有云笔记, 你可以下载它安装在自己的服务器上, 当然也可以在 http://gemsnote.com 上注册.\n\n这里详细整理了gemsnote二进版和gemsnote开发版的安装教程, 请移步至:\n\n* [gemsnote二进制详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n* [gemsnote开发版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n\n## 4. 如何对gemsnote进行二次开发\n\n请查看 [How-to-develop-gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91gemsnote)\n\n## 5 相关文档\n* [gemsnote二进制版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n* [gemsnote开发版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n* [Gemsnote source gemsnote源码导读](https://github.com/gemsnote/gemsnote/wiki/Gemsnote-source-gemsnote源码导读)\n* [gemsnote blog theme api(中文版)](https://github.com/gemsnote/gemsnote/wiki/gemsnote-blog-theme-api)\n* [How to develop gemsnote 如何开发gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-如何开发gemsnote)\n\n更多文档请查看 [wiki](https://github.com/gemsnote/gemsnote/wiki).\n\n## 6. 贡献者\n多谢 [贡献者](https://github.com/gemsnote/gemsnote/graphs/contributors) 的贡献, gemsnote因有你们而更完美!\n\n## 7. 加入我们\n\n欢迎提交[pull requests](https://github.com/gemsnote/gemsnote/pulls) 到gemsnote.\n\n有任何问题或建议, 欢迎提交[issue](https://github.com/gemsnote/gemsnote/issues).\n\nGemsnote还有很多问题, 如果你喜欢它, 欢迎加入我们一起完善gemsnote.\n\n## 8. 捐赠\n支持我们, [捐赠Gemsnote](http://gemsnote.org/#donate). 感谢[捐赠者](http://gemsnote.gemsnote.com/post/gemsnote-donation-list), 谢谢你们的鼓励, Gemsnote会一直坚持!\n\n## 9. 其它相关项目\n* [Gemsnote Desktop App](https://github.com/gemsnote/desktop-app), [下载地址](http://app.gemsnote.com)\n* [Gemsnote IOS](https://github.com/gemsnote/gemsnote-ios), 开发阶段\n* [Gemsnote Android](https://github.com/Dminter/gemsnote-android-client), 开发阶段\n\n同样, 欢迎加入我们!\n\n## 讨论\n* [gemsnote 社区](http://bbs.gemsnote.com)\n* QQ群: 158716820\n* [gemsnote google group](https://groups.google.com/forum/#!forum/gemsnote)\n\n----------------------------------------------------------------\n[English](README.md)", "UpdatedTime": "2015-06-15T18:36:35.265+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}, {"Content": "# Gemsnote产品\n\n## 1. 介绍\n\nGemsnote, 不只是笔记!\n![gemsnote](http://7xj51o.com1.z0.glb.clouddn.com/default_markdown.png)\n\n\n**特性**\n\n* 知识管理: 通过gemsnote来管理知识, gemsnote有易操作的界面, 包含两款编辑器tinymce和markdown. 在gemsnote, 你可以尽情享受写作.\n* 分享: 你也可以通过分享知识给好友, 让好友拥有你的知识.\n* 协作: 在分享的同时也可以与好友一起协作知识.\n* 博客: gemsnote也可以作为你的博客, 将知识公开成博客, 让gemsnote把你的知识传播的更远!\n\n## 2. 为什么我们要创建gemsnote?\n说实话, 我们曾是evernote的忠实粉丝, 但是我们也发现evernote的不足:\n* evernote的编辑器不能满足我们的需求, 不能贴代码(格式会乱掉, 作为程序员, 代码是我们的基本需求啊), 图片不能缩放.\n* 我们是markdown的爱好者, 可是evernote竟然没有.\n* 我们也想将知识公开, 所以我们有自己的博客, 如wordpress, 但为什么这两者不能合二为一呢?\n* 还有...\n\n## 3.安装gemsnote\ngemsnote是一款私有云笔记, 你可以下载它安装在自己的服务器上, 当然也可以在 http://gemsnote.com 上注册.\n\n这里详细整理了gemsnote二进版和gemsnote开发版的安装教程, 请移步至:\n\n* [gemsnote二进制详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n* [gemsnote开发版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n\n## 4. 如何对gemsnote进行二次开发\n\n请查看 [How-to-develop-gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91gemsnote)\n\n## 5 相关文档\n* [gemsnote二进制版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n* [gemsnote开发版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)\n* [Gemsnote source gemsnote源码导读](https://github.com/gemsnote/gemsnote/wiki/Gemsnote-source-gemsnote源码导读)\n* [gemsnote blog theme api(中文版)](https://github.com/gemsnote/gemsnote/wiki/gemsnote-blog-theme-api)\n* [How to develop gemsnote 如何开发gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-如何开发gemsnote)\n\n更多文档请查看 [wiki](https://github.com/gemsnote/gemsnote/wiki).\n\n## 6. 贡献者\n多谢 [贡献者](https://github.com/gemsnote/gemsnote/graphs/contributors) 的贡献, gemsnote因有你们而更完美!\n\n## 7. 加入我们\n\n欢迎提交[pull requests](https://github.com/gemsnote/gemsnote/pulls) 到gemsnote.\n\n有任何问题或建议, 欢迎提交[issue](https://github.com/gemsnote/gemsnote/issues).\n\nGemsnote还有很多问题, 如果你喜欢它, 欢迎加入我们一起完善gemsnote.\n\n## 8. 捐赠\n支持我们, [捐赠Gemsnote](http://gemsnote.org/#donate). 感谢[捐赠者](http://gemsnote.gemsnote.com/post/gemsnote-donation-list), 谢谢你们的鼓励, Gemsnote会一直坚持!\n\n## 9. 其它相关项目\n* [Gemsnote Desktop App](https://github.com/gemsnote/desktop-app), [下载地址](http://app.gemsnote.com)\n* [Gemsnote IOS](https://github.com/gemsnote/gemsnote-ios), 开发阶段\n* [Gemsnote Android](https://github.com/Dminter/gemsnote-android-client), 开发阶段\n\n同样, 欢迎加入我们!\n\n## 讨论\n* [gemsnote 社区](http://bbs.gemsnote.com)\n* QQ群: 158716820\n* [gemsnote google group](https://groups.google.com/forum/#!forum/gemsnote)\n\n----------------------------------------------------------------\n[English](README.md)", "UpdatedTime": "2015-06-15T18:36:34.059+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}]') ON CONFLICT DO NOTHING;
INSERT INTO public.note_content_histories (id, user_id, histories) VALUES ('561bb57505fcd164d3000000', '5368c1aa99c37b029d000001', '[{"Content": "<p>dddd</p>", "UpdatedTime": "2015-10-12T21:28:26.959+08:00", "UpdatedUserId": "5368c1aa99c37b029d000001"}]') ON CONFLICT DO NOTHING;


--
-- Data for Name: note_contents; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.note_contents (id, user_id, is_blog, content, abstract, created_time, updated_time, updated_user_id) VALUES ('5368c53199c37b095a000005', '5368c53199c37b095a000001', true, '<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:&nbsp;</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>gemsnote已开源, 欢迎加入我们, 让gemsnote更好!</p><p>github:&nbsp;<a href="https://github.com/gemsnote/gemsnote" data-mce-href="https://github.com/gemsnote/gemsnote">https://github.com/gemsnote/gemsnote</a></p><p>官网: <a href="http://gemsnote.com" data-mce-href="http://gemsnote.com">http://gemsnote.com</a></p>', '<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:&nbsp;</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>gemsnote已开源, 欢迎加入我们, 让gemsnote更好!</p><p>github:&nbsp;<a href="https://github.com/gemsnote/gemsnote" data-mce-href="https://github.com/gemsnote/gemsnote">https://github.com/gemsnote/gemsnote</a></p><p>官网: <a href="http://gemsnote.com" data-mce-href="http://gemsnote.com">http://gemsnote.com</a></p>', '2014-05-06 11:19:13.339+00', '2014-05-06 11:19:13.339+00', '5368c53199c37b095a000001') ON CONFLICT DO NOTHING;
INSERT INTO public.note_contents (id, user_id, is_blog, content, abstract, created_time, updated_time, updated_user_id) VALUES ('540814f199c37b555d000005', '540814f199c37b555d000001', true, '<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:&nbsp;</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>gemsnote已开源, 欢迎加入我们, 让gemsnote更好!</p><p>github:&nbsp;<a href="https://github.com/gemsnote/gemsnote" data-mce-href="https://github.com/gemsnote/gemsnote">https://github.com/gemsnote/gemsnote</a></p><p>官网: <a href="http://gemsnote.com" data-mce-href="http://gemsnote.com">http://gemsnote.com</a></p>', '<p>Hi, 你好!&nbsp;欢迎来到gemsnote!</p><p>在这里, 你可以:&nbsp;</p><ul><li>管理知识</li><li>分享知识</li><li>与好友协作知识</li><li>将知识公开成博客</li><li>.........等待你来发现, 等待我们为你呈现更好的gemsnote!<br></li></ul><p>gemsnote已开源, 欢迎加入我们, 让gemsnote更好!</p><p>github:&nbsp;<a href="https://github.com/gemsnote/gemsnote" data-mce-href="https://github.com/gemsnote/gemsnote">https://github.com/gemsnote/gemsnote</a></p><p>官网: <a href="http://gemsnote.com" data-mce-href="http://gemsnote.com">http://gemsnote.com</a></p>', '2014-09-04 07:29:53.192+00', '2014-09-04 07:29:53.192+00', '540814f199c37b555d000001') ON CONFLICT DO NOTHING;
INSERT INTO public.note_contents (id, user_id, is_blog, content, abstract, created_time, updated_time, updated_user_id) VALUES ('5483207cf4e87203a4000001', '5368c1aa99c37b029d000001', false, '
# About Gemsnote

## 1. Introduction

Gemsnote, not just a notepad!
![gemsnote](http://7xj51o.com1.z0.glb.clouddn.com/default_markdown.png)

**Some Features**

* Knowledge: Manage your knowledge in gemsnote. gemsnote contains the tinymce editor and a markdown editor, just enjoy yourself writing.
* Share: Share your knowledge with your friends in gemsnote. You can invite your friends to join your notepad in the cloud so you can share knowledge.
* Cooperation: Collaborate with friends to improve your skills.
* Blog: Publish your knowledge and make gemsnote your blog.

## 2. Why we created gemsnote
To be honest, our inspiration comes from Evernote. We use Evernote to manage our knowledge everyday. But we find that:
* Evernote''s editor can''t meet our needs, it does not have document navigation, it does not render code properly (as a programmer, syntax highlighted code rendering is a basic need), it cannot resize images and so forth
* We like markdown, but Evernote does not support it.
* We want to share our knowledge, so all of us have our blogs (e.g. on Wordpress) and our Evernote accounts, but why can not those two be one!
* ......

## 3. How to install gemsnote

More information about how to install gemsnote please see:
* [gemsnote binary distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-binary-distribution-installation-tutorial)
* [gemsnote develop distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-develop-distribution-installation-tutorial)

## 4. How to develop gemsnote

Please see [How-to-develop-gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91gemsnote)


## 5. Docs
* [gemsnote binary distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-binary-distribution-installation-tutorial)
* [gemsnote develop distribution installation tutorial](https://github.com/gemsnote/gemsnote/wiki/gemsnote-develop-distribution-installation-tutorial)
* [gemsnote blog theme api](https://github.com/gemsnote/gemsnote/wiki/gemsnote-blog-theme-api_en)

More docs please see [wiki](https://github.com/gemsnote/gemsnote/wiki).

## 6. Contributors
Thank you to all the [contributors](https://github.com/gemsnote/gemsnote/graphs/contributors) on
this project. Your help is much appreciated.

## 7.Join us
Please fork this repository and contribute back using [pull requests](https://github.com/gemsnote/gemsnote/pulls).

If you find some problems or has some good ideas, please submit [issues](https://github.com/gemsnote/gemsnote/issues).

You are always welcomed!

## 8. Donation
Support us, [donate us](http://gemsnote.org/#donate). And thanks [donators](http://gemsnote.gemsnote.com/post/gemsnote-donation-list).

## 9. Related projects
* [Gemsnote Desktop App](https://github.com/gemsnote/desktop-app), [Download](http://app.gemsnote.com)
* [Gemsnote IOS](https://github.com/gemsnote/gemsnote-ios), development phase
* [Gemsnote Android](https://github.com/Dminter/gemsnote-android-client), development phase

And also, you are welcome to join us.

## 9. Discussion
* [gemsnote bbs](http://bbs.gemsnote.com)
* [gemsnote google group](https://groups.google.com/forum/#!forum/gemsnote)
* QQ Group: 158716820


', '<div>
						                        <div id="wmd-preview" class="preview-content"></div>
						                    <div id="wmd-preview-section-500" class="wmd-preview-section preview-content">

</div><div id="wmd-preview-section-501" class="wmd-preview-section preview-content">

<h1 id="about-gemsnote">About Gemsnote</h1>

</div><div id="wmd-preview-section-502" class="wmd-preview-section preview-content">

<h2 id="1-introduction">1. Introduction</h2>

<p>Gemsnote, not just a notepad! <br>
<img src="http://7xj51o.com1.z0.glb.clouddn.com/default_markdown.png" alt="gemsnote" title=""></p>

<p><strong>Some Features</strong></p>

<ul>
<li>Knowledge: Manage your knowledge in gemsnote. gemsnote contains the tinymce editor and a markdown editor, just enjoy yourself writing.</li>
<li>Share: Share your knowledge with your friends in gemsnote. You can invite your friends to join your notepad in the cloud so you can share knowledge.</li>
<li>Cooperation: Collaborate with friends to improve your skills.</li>
<li>Blog: Publish your knowledge and make gemsnote your blog.</li>
</ul>

</div><div id="wmd-preview-section-503" class="wmd-preview-section preview-content">

<h2 id="2-why-we-created-gemsnote">2. Why we created gemsnote</h2>

<p>To be honest, our inspiration comes from Evernote. We use Evernote to manage our knowledge everyday. But we find that: <br>
* Evernote’s editor can’t meet our needs, it does not have document navigation, it does not render code properly (as a programmer, syntax highlighted code rendering is a basic need), it cannot resize images and so forth <br>
* We like markdown, but Evernote does not support it. <br>
* We want to share our knowle</p></div></div>', '2014-12-06 15:28:01.82+00', '2015-06-15 10:37:30.723+00', '5368c1aa99c37b029d000001') ON CONFLICT DO NOTHING;
INSERT INTO public.note_contents (id, user_id, is_blog, content, abstract, created_time, updated_time, updated_user_id) VALUES ('5524ba2f99c37b292000000b', '5524ba2f99c37b2920000007', false, '# Gemsnote产品说明

## 1. 介绍

Gemsnote, 不只是笔记!

**特性**

* 知识管理: 通过gemsnote来管理知识, gemsnote有易操作的界面, 包含两款编辑器tinymce和markdown. 在gemsnote, 你可以尽情享受写作.
* 分享: 你也可以通过分享知识给好友, 让好友拥有你的知识.
* 协作: 在分享的同时也可以与好友一起协作知识.
* 博客: gemsnote也可以作为你的博客, 将知识公开成博客, 让gemsnote把你的知识传播的更远!

## 2. 为什么我们要创建gemsnote?
说实话, 我们曾是evernote的忠实粉丝, 但是我们也发现evernote的不足:
* evernote的编辑器不能满足我们的需求, 不能贴代码(格式会乱掉, 作为程序员, 代码是我们的基本需求啊), 图片不能缩放.
* 我们是markdown的爱好者, 可是evernote竟然没有.
* 我们也想将知识公开, 所以我们有自己的博客, 如wordpress, 但为什么这两者不能合二为一呢?
* 还有...

## 3.安装gemsnote
gemsnote是一款私有云笔记, 你可以下载它安装在自己的服务器上, 当然也可以在 http://gemsnote.com 上注册.

这里详细整理了gemsnote二进版和gemsnote开发版的安装教程, 请移步至:
* [gemsnote二进制详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)
* [gemsnote开发版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)

### 3.1. 下载gemsnote

Gemsnote V1.0-beta 已发布, 二进制文件(暂时没有windows版的):

* Linux: [gemsnote-linux-x86_64.v1.0-beta.bin.tar.gz](https://github.com/gemsnote/gemsnote/releases/download/1.0-beta/gemsnote-linux-x86_64.v1.0-beta.bin.tar.gz)
* MacOS X: [gemsnote-mac-x86_64.v1.0-beta.bin.tar.gz](https://github.com/gemsnote/gemsnote/releases/download/1.0-beta/gemsnote-mac-x86_64.v1.0-beta.bin.tar.gz)

### 3.2. 安装 MongodbDB

Gemsnote是由golang(使用[revel](https://revel.github.io/)框架 和 [MongoDB](https://www.mongodb.org)数据库), 你需要先安装Mongodb.

安装MongodbDB, 导入数据更多细节请查看: [wiki](https://github.com/gemsnote/gemsnote/wiki/Install-Mongodb)

### 3.3. 导入初始数据

MongodbDB初始数据在 `[PATH_TO_GEMSNOTE]/mongodb_backup/gemsnote_install_data`

```
$> mongorestore -h localhost -d gemsnote --directoryperdb PATH_TO_GEMSNOTE/mongodb_backup/gemsnote_install_data
```

初始数据包含两个用户:

```
user2 username: admin, password: abc123 (管理员, 重要!)
user3 username: demo@gemsnote.com, password: demo@gemsnote.com (为体验使用)
```

### 3.4. 配置

修改 `[PATH_TO_GEMSNOTE]/conf/app.conf`. 有以下选项:

``mongodb``  **必须配置!**

```Shell
db.host=localhost
db.port=27017
db.dbname=gemsnote
db.username=
db.password=
```

``app.secret`` **重要**
请随意修改一个, app的密钥, 不能使用默认的, 不然会有安全问题

更多配置请查看 `app/app.conf` 和 [revel 手册](https://revel.github.io/)

### 3.5. 运行gemsnote

```
$> cd PATH_TO_GEMSNOTE/bin
$> sudo sh run.sh
```

## 4. 如何对gemsnote进行二次开发

请查看 [How-to-develop-gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91gemsnote)

## 5. 贡献者
多谢 [贡献者](https://github.com/gemsnote/gemsnote/graphs/contributors) 的贡献, gemsnote因有你们而更完美!

## 6. 加入我们

欢迎提交[pull requests](https://github.com/gemsnote/gemsnote/pulls) 到gemsnote.

gemsnote还有很多问题, 如果你喜欢它, 欢迎加入我们一起完善gemsnote.

## 讨论
* [gemsnote 社区](http://bbs.gemsnote.com)
* QQ群: 158716820
* [gemsnote google group](https://groups.google.com/forum/#!forum/gemsnote)', '<h1>Gemsnote产品说明</h1>

<h2>1. 介绍</h2>

<p>Gemsnote, 不只是笔记!</p>

<p><strong>特性</strong></p>

<ul>
<li>知识管理: 通过gemsnote来管理知识, gemsnote有易操作的界面, 包含两款编辑器tinymce和markdown. 在gemsnote, 你可以尽情享受写作.</li>
<li>分享: 你也可以通过分享知识给好友, 让好友拥有你的知识.</li>
<li>协作: 在分享的同时也可以与好友一起协作知识.</li>
<li>博客: gemsnote也可以作为你的博客, 将知识公开成博客, 让gemsnote把你的知识传播的更远!</li>
</ul>

<h2>2. 为什么我们要创建gemsnote?</h2>

<p>说实话, 我们曾是evernote的忠实粉丝, 但是我们也发现evernote的不足:
* evernote的编辑器不能满足我们的需求, 不能贴代码(格式会乱掉, 作为程序员, 代码是我们的基本需求啊), 图片不能缩放.
* 我们是markdown的爱好者, 可是evernote竟然没有.
* 我们也想将知识公开, 所以我们有自己的博客, 如wordpress, 但为什么这两者不能合二为一呢?
* 还有...</p>

<h2>3.安装gemsnote</h2>

<p>gemsnote是一款私有云笔记, 你可以下载它安装在自己的服务器上, 当然也可以在 <a href="http://gemsnote.com">http://gemsnote.com</a> 上注册.</p>

<p>这里详细整理了gemsnote二进版和gemsnote开发版的安装教程, 请移步至:
* <a href="https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B">gemsnote二进制详细安装教程</a>
* <a href="https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B">gemsnote开发版详细安装教程</a></p>

<h3>3.1. 下载gemsnote</h3>

<p>Gemsnote V1.0-beta 已发布, 二进制文件(暂时没有windows版的):</p>

<ul>
<li>Linux: <a href="https://github.com/gemsnote/gemsnote/releases/download/1.0-beta/gemsnote-linux-x86_64.v1.0-beta.bin.tar.gz">gemsnote-linux-x86_64.v1.0-beta.bin.tar.gz</a></li>
<li>MacOS X: <a href="https://github.com/gemsnote/gemsnote/releases/download/1.0-beta/gemsnote-mac-x86_64.v1.0-beta.bin.tar.gz">gemsnote-mac-x86_64.v1.0-beta.bin.tar.gz</a></li>
</ul>

<h3>3.2. 安装 MongodbDB</h3>

<p>Gemsnote是由golang(使用<a href="https://revel.github.io/">revel</a>框架 和 <a href="https://www.mongodb.org">MongoDB</a>数据库), 你需要先安装Mongodb.</p>

<p>安装MongodbDB, 导入数据更多细节请查看: <a href="https://github.com/gemsnote/gemsnote/wiki/Install-Mongodb">wiki</a></p>

<h3>3.3. 导入初始数据</h3>

<p>MongodbDB初始数据在 <code>[PATH_TO_GEMSNOTE]/mongodb_backup/gemsnote_install_data</code></p>

<pre class="prettyprint linenums prettyprinted"><ol class="linenums"><li class="L0"><code><span class="pln">$</span><span class="pun">&gt;</span><span class="pln"> mongorestore </span><span class="pun">-</span><span class="pln">h loca</span></code></li></ol></pre>', '2015-04-08 05:18:39.433+00', '2015-04-08 05:18:39.433+00', '5524ba2f99c37b2920000007') ON CONFLICT DO NOTHING;
INSERT INTO public.note_contents (id, user_id, is_blog, content, abstract, created_time, updated_time, updated_user_id) VALUES ('557eaa9905fcd14d95000001', '5368c1aa99c37b029d000001', false, '# 关于Gemsnote

## 1. 介绍

Gemsnote, 不只是笔记!
![gemsnote](http://7xj51o.com1.z0.glb.clouddn.com/default_markdown.png)

**特性**

* 知识管理: 通过gemsnote来管理知识, gemsnote有易操作的界面, 包含两款编辑器tinymce和markdown. 在gemsnote, 你可以尽情享受写作.
* 分享: 你也可以通过分享知识给好友, 让好友拥有你的知识.
* 协作: 在分享的同时也可以与好友一起协作知识.
* 博客: gemsnote也可以作为你的博客, 将知识公开成博客, 让gemsnote把你的知识传播的更远!

## 2. 为什么我们要创建gemsnote?
说实话, 我们曾是evernote的忠实粉丝, 但是我们也发现evernote的不足:
* evernote的编辑器不能满足我们的需求, 不能贴代码(格式会乱掉, 作为程序员, 代码是我们的基本需求啊), 图片不能缩放.
* 我们是markdown的爱好者, 可是evernote竟然没有.
* 我们也想将知识公开, 所以我们有自己的博客, 如wordpress, 但为什么这两者不能合二为一呢?
* 还有...

## 3.安装gemsnote
gemsnote是一款私有云笔记, 你可以下载它安装在自己的服务器上, 当然也可以在 http://gemsnote.com 上注册.

这里详细整理了gemsnote二进版和gemsnote开发版的安装教程, 请移步至:

* [gemsnote二进制详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)
* [gemsnote开发版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)

## 4. 如何对gemsnote进行二次开发

请查看 [How-to-develop-gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91gemsnote)

## 5 相关文档
* [gemsnote二进制版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)
* [gemsnote开发版详细安装教程](https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)
* [Gemsnote source gemsnote源码导读](https://github.com/gemsnote/gemsnote/wiki/Gemsnote-source-gemsnote源码导读)
* [gemsnote blog theme api(中文版)](https://github.com/gemsnote/gemsnote/wiki/gemsnote-blog-theme-api)
* [How to develop gemsnote 如何开发gemsnote](https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-如何开发gemsnote)

更多文档请查看 [wiki](https://github.com/gemsnote/gemsnote/wiki).

## 6. 贡献者
多谢 [贡献者](https://github.com/gemsnote/gemsnote/graphs/contributors) 的贡献, gemsnote因有你们而更完美!

## 7. 加入我们

欢迎提交[pull requests](https://github.com/gemsnote/gemsnote/pulls) 到gemsnote.

有任何问题或建议, 欢迎提交[issue](https://github.com/gemsnote/gemsnote/issues).

Gemsnote还有很多问题, 如果你喜欢它, 欢迎加入我们一起完善gemsnote.

## 8. 捐赠
支持我们, [捐赠Gemsnote](http://gemsnote.org/#donate). 感谢[捐赠者](http://gemsnote.gemsnote.com/post/gemsnote-donation-list), 谢谢你们的鼓励, Gemsnote会一直坚持!

## 9. 其它相关项目
* [Gemsnote Desktop App](https://github.com/gemsnote/desktop-app), [下载地址](http://app.gemsnote.com)
* [Gemsnote IOS](https://github.com/gemsnote/gemsnote-ios), 开发阶段
* [Gemsnote Android](https://github.com/Dminter/gemsnote-android-client), 开发阶段

同样, 欢迎加入我们!

## 讨论
* [gemsnote 社区](http://bbs.gemsnote.com)
* QQ群: 158716820
* [gemsnote google group](https://groups.google.com/forum/#!forum/gemsnote)', '<div>
						                        <div id="wmd-preview" class="preview-content"></div>
						                    <div id="wmd-preview-section-452" class="wmd-preview-section preview-content">

</div><div id="wmd-preview-section-453" class="wmd-preview-section preview-content">

<h1 id="关于gemsnote">关于Gemsnote</h1>

</div><div id="wmd-preview-section-454" class="wmd-preview-section preview-content">

<h2 id="1-介绍">1. 介绍</h2>

<p>Gemsnote, 不只是笔记! <br>
<img src="http://7xj51o.com1.z0.glb.clouddn.com/default_markdown.png" alt="gemsnote" title=""></p>

<p><strong>特性</strong></p>

<ul>
<li>知识管理: 通过gemsnote来管理知识, gemsnote有易操作的界面, 包含两款编辑器tinymce和markdown. 在gemsnote, 你可以尽情享受写作.</li>
<li>分享: 你也可以通过分享知识给好友, 让好友拥有你的知识.</li>
<li>协作: 在分享的同时也可以与好友一起协作知识.</li>
<li>博客: gemsnote也可以作为你的博客, 将知识公开成博客, 让gemsnote把你的知识传播的更远!</li>
</ul>

</div><div id="wmd-preview-section-455" class="wmd-preview-section preview-content">

<h2 id="2-为什么我们要创建gemsnote">2. 为什么我们要创建gemsnote?</h2>

<p>说实话, 我们曾是evernote的忠实粉丝, 但是我们也发现evernote的不足: <br>
* evernote的编辑器不能满足我们的需求, 不能贴代码(格式会乱掉, 作为程序员, 代码是我们的基本需求啊), 图片不能缩放. <br>
* 我们是markdown的爱好者, 可是evernote竟然没有. <br>
* 我们也想将知识公开, 所以我们有自己的博客, 如wordpress, 但为什么这两者不能合二为一呢? <br>
* 还有…</p>

</div><div id="wmd-preview-section-456" class="wmd-preview-section preview-content">

<h2 id="3安装gemsnote">3.安装gemsnote</h2>

<p>gemsnote是一款私有云笔记, 你可以下载它安装在自己的服务器上, 当然也可以在 <a href="http://gemsnote.com">http://gemsnote.com</a> 上注册.</p>

<p>这里详细整理了gemsnote二进版和gemsnote开发版的安装教程, 请移步至:</p>

<ul>
<li><a href="https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B">gemsnote二进制详细安装教程</a></li>
<li><a href="https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B">gemsnote开发版详细安装教程</a></li>
</ul>

</div><div id="wmd-preview-section-457" class="wmd-preview-section preview-content">

<h2 id="4-如何对gemsnote进行二次开发">4. 如何对gemsnote进行二次开发</h2>

<p>请查看 <a href="https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91gemsnote">How-to-develop-gemsnote</a></p>

</div><div id="wmd-preview-section-458" class="wmd-preview-section preview-content">

<h2 id="5-相关文档">5 相关文档</h2>

<ul>
<li><a href="https://github.com/gemsnote/gemsnote/wiki/gemsnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B">gemsnote二进制版详细安装教程</a></li>
<li><a href="https://github.com/gemsnote/gemsnote/wiki/gemsnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B">gemsnote开发版详细安装教程</a></li>
<li><a href="https://github.com/gemsnote/gemsnote/wiki/Gemsnote-source-gemsnote源码导读">Gemsnote source gemsnote源码导读</a></li>
<li><a href="https://github.com/gemsnote/gemsnote/wiki/gemsnote-blog-theme-api">gemsnote blog theme api(中文版)</a></li>
<li><a href="https://github.com/gemsnote/gemsnote/wiki/How-to-develop-gemsnote-如何开发gemsnote">How to develop gemsnote 如何开发gemsnote</a></li>
</ul>

<p>更多文档请查看 <a href="https://github.com/gemsnote/gemsnote/wiki">wiki</a>.</p>

</div><div id="wmd-preview-section-459" class="wmd-preview-section preview-content">

<h2 id="6-贡献者">6. 贡献者</h2>

<p>多谢 <a href="https://github.com/gemsnote/gemsnote/graphs/contributors">贡献者</a> 的贡献, gemsnote因有你们而更完美!</p>

</div><div id="wmd-preview-section-460" class="wmd-preview-section preview-content">

<h2 id="7-加入我们">7. 加入我们</h2>

<p>欢迎提交<a href="https://github.com/gemsnote/gemsnote/pulls">pull requests</a> 到gemsnote.</p>

<p>有任何问题或建议, 欢迎提交<a href="https://github.com/gemsnote/gemsnote/issues">issue</a>.</p>

<p>Le</p></div></div>', '2015-06-15 10:36:13.335+00', '2015-06-15 10:37:24.149+00', '5368c1aa99c37b029d000001') ON CONFLICT DO NOTHING;


--
-- Data for Name: note_images; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541bf730c1af05c13c805401', '5412e2a219807a68c3000000', '541bf1ca99c37b7f0a000002') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541bf730c1af05c13c805402', '5412e2a219807a68c3000000', '541bf21099c37b7f0a000004') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541bf730c1af05c13c805403', '5412e2a219807a68c3000000', '541bf36b99c37b7f0a000007') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541bf730c1af05c13c805404', '5412e2a219807a68c3000000', '541bf45499c37b824d000002') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541bf730c1af05c13c805405', '5412e2a219807a68c3000000', '541bf1af99c37b7f0a000001') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541bf730c1af05c13c805406', '5412e2a219807a68c3000000', '541bf64399c37b824d000004') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541bf730c1af05c13c805407', '5412e2a219807a68c3000000', '541bf68099c37b824d000008') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541bf730c1af05c13c805408', '5412e2a219807a68c3000000', '541bf70699c37b824d00000c') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541d253c2ddddba10064c340', '541d24c209915b0b5d000000', '541d24d499c37b0eb500000a') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541d255e2ddddba10064c342', '541d253c09915b0b5d000001', '541d254999c37b0eb500000d') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541e37300716dcbd8f1f647a', '541d23c399c37b0eb5000001', '541e35cc99c37b06ae000002') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541e37300716dcbd8f1f647b', '541d23c399c37b0eb5000001', '541d23c399c37b0eb5000002') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541e37300716dcbd8f1f647c', '541d23c399c37b0eb5000001', '541d23c399c37b0eb5000003') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541e37300716dcbd8f1f647d', '541d23c399c37b0eb5000001', '541d23c399c37b0eb5000004') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541e37300716dcbd8f1f647e', '541d23c399c37b0eb5000001', '541d23c399c37b0eb5000005') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541e37300716dcbd8f1f647f', '541d23c399c37b0eb5000001', '541d23c399c37b0eb5000006') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541e37300716dcbd8f1f6480', '541d23c399c37b0eb5000001', '541d23c399c37b0eb5000007') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541e37300716dcbd8f1f6481', '541d23c399c37b0eb5000001', '541d23c399c37b0eb5000008') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541e37300716dcbd8f1f6482', '541d23c399c37b0eb5000001', '541d23c399c37b0eb5000009') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541eb97f0716dcbd8f1f6497', '541d25dc09915b19a9000000', '541d25ec99c37b11be000001') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541eb97f0716dcbd8f1f6498', '541d25dc09915b19a9000000', '541daa9499c37b0ef5000004') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541eb97f0716dcbd8f1f6499', '541d25dc09915b19a9000000', '541daabd99c37b0ef5000006') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('541eb97f0716dcbd8f1f649a', '541d25dc09915b19a9000000', '541daad399c37b0ef5000008') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('5425213a3d1a833d97d7ae85', '541d2a7f99c37b1947000001', '541e31ed99c37b06ae000001') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('5425213a3d1a833d97d7ae86', '541d2a7f99c37b1947000001', '54250c1a99c37b37f6000001') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('5425213a3d1a833d97d7ae87', '541d2a7f99c37b1947000001', '54251fe299c37b5517000001') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('5425213a3d1a833d97d7ae88', '541d2a7f99c37b1947000001', '541edd8b99c37be24a000001') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('5425213a3d1a833d97d7ae89', '541d2a7f99c37b1947000001', '541daadc99c37b0ef5000009') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('5425213a3d1a833d97d7ae8a', '541d2a7f99c37b1947000001', '541daad399c37b0ef5000007') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('5425213a3d1a833d97d7ae8b', '541d2a7f99c37b1947000001', '541d256299c37b0eb500000f') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('5425213a3d1a833d97d7ae8c', '541d2a7f99c37b1947000001', '541d2c6c99c37b1980000002') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('54256b8b3d1a833d97d7aedc', '5371aa4e19807a273a000000', '5425678c99c37bb31a000001') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('54256b8b3d1a833d97d7aedd', '5371aa4e19807a273a000000', '541bf68099c37b824d000007') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('54256b8b3d1a833d97d7aede', '5371aa4e19807a273a000000', '54256a5e99c37bb31a000003') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('54256b8b3d1a833d97d7aedf', '5371aa4e19807a273a000000', '54256a1799c37bb31a000002') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('54256b8b3d1a833d97d7aee0', '5371aa4e19807a273a000000', '54256a9599c37bb31a000004') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('54256b8b3d1a833d97d7aee1', '5371aa4e19807a273a000000', '54252a0899c37b6453000001') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('54257d3f3d1a833d97d7aee2', '541bf33019807a69be000000', '54252a0899c37b6453000001') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('54257d3f3d1a833d97d7aee3', '541bf33019807a69be000000', '541bf64399c37b824d000003') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('54257d3f3d1a833d97d7aee4', '541bf33019807a69be000000', '541bf70699c37b824d00000b') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('54257d3f3d1a833d97d7aee5', '541bf33019807a69be000000', '541bf6b899c37b824d000009') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('54257d3f3d1a833d97d7aee6', '541bf33019807a69be000000', '541bf68099c37b824d000007') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('54257d3f3d1a833d97d7aee7', '541bf33019807a69be000000', '541bf20f99c37b7f0a000003') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('54257d3f3d1a833d97d7aee8', '541bf33019807a69be000000', '541bf33b99c37b7f0a000005') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('54257d3f3d1a833d97d7aee9', '541bf33019807a69be000000', '541bf36a99c37b7f0a000006') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('54257d3f3d1a833d97d7aeea', '541bf33019807a69be000000', '541bf45299c37b824d000001') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('542970ab899df964c82f033d', '542966f499c37bc034000007', '5429676a99c37bc034000008') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('542970e0899df964c82f033e', '542970d9e527672881000000', '5429676a99c37bc034000008') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('542a64a6ddc557eb8875ac28', '54267aa3e527673016000000', '5429646599c37bbfb0000001') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('542aa2a5ddc557eb8875ac2f', '542a1c59e527674305000000', '542966d999c37bc034000002') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('542aab58ddc557eb8875ac30', '541eb10019807a3063000000', '541efb9899c37b029600000d') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('5432ccb65f54348e5a17eeb6', '542003f319807a0ca8000000', '541e31ed99c37b06ae000001') ON CONFLICT DO NOTHING;
INSERT INTO public.note_images (id, note_id, image_id) VALUES ('5432ccb75f54348e5a17eeb7', '541eb0e219807a0bfc000000', '541e31ed99c37b06ae000001') ON CONFLICT DO NOTHING;


--
-- Data for Name: note_tags; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.note_tags (id, user_id, tag, usn, count, created_time, updated_time, is_deleted) VALUES ('551a3f5c99c37b04de000001', '5368c1aa99c37b029d000001', 'gemsnote', 1, 2, '2015-03-31 06:31:56.604+00', '2015-06-15 10:46:41.958+00', false) ON CONFLICT DO NOTHING;
INSERT INTO public.note_tags (id, user_id, tag, usn, count, created_time, updated_time, is_deleted) VALUES ('551a3f5c99c37b04de000002', '5368c1aa99c37b029d000001', '欢迎', 2, 1, '2015-03-31 06:31:56.604+00', '2015-03-31 06:31:56.604+00', false) ON CONFLICT DO NOTHING;
INSERT INTO public.note_tags (id, user_id, tag, usn, count, created_time, updated_time, is_deleted) VALUES ('551a3f5c99c37b04de000003', '5368c1aa99c37b029d000001', 'red', 3, 1, '2015-03-31 06:31:56.604+00', '2015-03-31 06:31:56.604+00', false) ON CONFLICT DO NOTHING;
INSERT INTO public.note_tags (id, user_id, tag, usn, count, created_time, updated_time, is_deleted) VALUES ('56238fb899c37b5318000001', '5368c1aa99c37b029d000001', '', 200035, 0, '2015-10-18 12:25:28.29+00', '2015-10-18 12:25:30.156+00', false) ON CONFLICT DO NOTHING;


--
-- Data for Name: notebooks; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.notebooks (id, user_id, parent_notebook_id, seq, title, url_title, number_notes, is_trash, is_blog, created_time, updated_time, usn, is_deleted) VALUES ('5368c9fc99c37b095a000007', '5368c9fc99c37b095a000006', NULL, -1, 'life', 'life', 0, false, false, '2026-09-01 03:31:54.960496+00', '2026-09-01 03:31:54.960496+00', 2, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notebooks (id, user_id, parent_notebook_id, seq, title, url_title, number_notes, is_trash, is_blog, created_time, updated_time, usn, is_deleted) VALUES ('5368c9fc99c37b095a000008', '5368c9fc99c37b095a000006', NULL, -1, 'study', 'study', 0, false, false, '2026-09-01 03:31:54.963+00', '2026-09-01 03:31:54.963+00', 3, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notebooks (id, user_id, parent_notebook_id, seq, title, url_title, number_notes, is_trash, is_blog, created_time, updated_time, usn, is_deleted) VALUES ('5368c9fc99c37b095a000009', '5368c9fc99c37b095a000006', NULL, -1, 'work', 'work', 0, false, false, '2026-09-01 03:31:54.965107+00', '2026-09-01 03:31:54.965107+00', 4, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notebooks (id, user_id, parent_notebook_id, seq, title, url_title, number_notes, is_trash, is_blog, created_time, updated_time, usn, is_deleted) VALUES ('540814f199c37b555d000002', '540814f199c37b555d000001', NULL, -1, 'life', 'life', 0, false, false, '2026-09-01 03:31:54.967198+00', '2026-09-01 03:31:54.967198+00', 5, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notebooks (id, user_id, parent_notebook_id, seq, title, url_title, number_notes, is_trash, is_blog, created_time, updated_time, usn, is_deleted) VALUES ('540814f199c37b555d000003', '540814f199c37b555d000001', NULL, -1, 'study', 'study', 0, false, false, '2026-09-01 03:31:54.968864+00', '2026-09-01 03:31:54.968864+00', 6, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notebooks (id, user_id, parent_notebook_id, seq, title, url_title, number_notes, is_trash, is_blog, created_time, updated_time, usn, is_deleted) VALUES ('540814f199c37b555d000004', '540814f199c37b555d000001', NULL, -1, 'work', 'work', 0, false, false, '2026-09-01 03:31:54.970679+00', '2026-09-01 03:31:54.970679+00', 7, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notebooks (id, user_id, parent_notebook_id, seq, title, url_title, number_notes, is_trash, is_blog, created_time, updated_time, usn, is_deleted) VALUES ('548125adf4e872105c000007', '5368c1aa99c37b029d000001', NULL, 1, 'Life', 'note', 2, false, false, '2026-09-01 03:31:54.978529+00', '2026-09-01 03:31:54.978529+00', 200024, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notebooks (id, user_id, parent_notebook_id, seq, title, url_title, number_notes, is_trash, is_blog, created_time, updated_time, usn, is_deleted) VALUES ('5524b99b99c37b2920000003', '5524b99b99c37b2920000002', NULL, -1, 'life', 'life', 0, false, false, '2015-04-08 05:16:11.704+00', '2015-04-08 05:16:11.704+00', 2, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notebooks (id, user_id, parent_notebook_id, seq, title, url_title, number_notes, is_trash, is_blog, created_time, updated_time, usn, is_deleted) VALUES ('5524b99b99c37b2920000004', '5524b99b99c37b2920000002', NULL, -1, 'study', 'study', 0, false, false, '2015-04-08 05:16:11.705+00', '2015-04-08 05:16:11.705+00', 3, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notebooks (id, user_id, parent_notebook_id, seq, title, url_title, number_notes, is_trash, is_blog, created_time, updated_time, usn, is_deleted) VALUES ('5524b99b99c37b2920000005', '5524b99b99c37b2920000002', NULL, -1, 'work', 'work', 0, false, false, '2015-04-08 05:16:11.7+00', '2015-04-08 05:16:11.7+00', 1, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notebooks (id, user_id, parent_notebook_id, seq, title, url_title, number_notes, is_trash, is_blog, created_time, updated_time, usn, is_deleted) VALUES ('5524ba2f99c37b2920000008', '5524ba2f99c37b2920000007', NULL, -1, 'life', 'life', 1, false, false, '2015-04-08 05:18:39.422+00', '2015-04-08 05:18:39.422+00', 3, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notebooks (id, user_id, parent_notebook_id, seq, title, url_title, number_notes, is_trash, is_blog, created_time, updated_time, usn, is_deleted) VALUES ('5524ba2f99c37b2920000009', '5524ba2f99c37b2920000007', NULL, -1, 'study', 'study', 0, false, false, '2015-04-08 05:18:39.419+00', '2015-04-08 05:18:39.419+00', 1, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notebooks (id, user_id, parent_notebook_id, seq, title, url_title, number_notes, is_trash, is_blog, created_time, updated_time, usn, is_deleted) VALUES ('5524ba2f99c37b292000000a', '5524ba2f99c37b2920000007', NULL, -1, 'work', 'work', 0, false, false, '2015-04-08 05:18:39.42+00', '2015-04-08 05:18:39.42+00', 2, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notebooks (id, user_id, parent_notebook_id, seq, title, url_title, number_notes, is_trash, is_blog, created_time, updated_time, usn, is_deleted) VALUES ('557eab5705fcd14d95000002', '5368c1aa99c37b029d000001', NULL, 0, 'Work', 'Work', 0, false, false, '2015-06-15 10:39:21.399+00', '2015-06-15 10:39:21.399+00', 200023, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notebooks (id, user_id, parent_notebook_id, seq, title, url_title, number_notes, is_trash, is_blog, created_time, updated_time, usn, is_deleted) VALUES ('557eab5c05fcd14d95000003', '5368c1aa99c37b029d000001', NULL, 2, 'Others', 'Others', 0, false, false, '2015-06-15 10:39:26.98+00', '2015-06-15 10:39:26.98+00', 200025, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notebooks (id, user_id, parent_notebook_id, seq, title, url_title, number_notes, is_trash, is_blog, created_time, updated_time, usn, is_deleted) VALUES ('557eab6705fcd14d95000004', '5368c1aa99c37b029d000001', '557eab5c05fcd14d95000003', -1, 'Travel', 'Untitled-5368c1aa99c37b029d000001', 0, false, false, '2015-06-15 10:39:36.345+00', '2015-06-15 10:39:36.345+00', 200027, false) ON CONFLICT DO NOTHING;


--
-- Data for Name: notes; Type: TABLE DATA; Schema: public;
INSERT INTO public.notes (id, user_id, created_user_id, notebook_id, title, description, src, img_src, tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined, read_num, like_num, comment_num, is_markdown, attach_num, created_time, updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted) VALUES ('5481481bf4e87273d2000003', '5368c1aa99c37b029d000001', NULL, '54814817f4e87273d2000002', '都放到', '撒旦法 ', '', '', '{}', true, false, '%E9%83%BD%E6%94%BE%E5%88%B0', false, false, false, 0, 0, 0, false, 0, '2014-12-05 05:52:32.856+00', '2014-12-05 05:52:32.856+00', NULL, NULL, '5368c1aa99c37b029d000001', 200044, true) ON CONFLICT DO NOTHING;
INSERT INTO public.notes (id, user_id, created_user_id, notebook_id, title, description, src, img_src, tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined, read_num, like_num, comment_num, is_markdown, attach_num, created_time, updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted) VALUES ('5481489cf4e8721ee3000000', '5368c1aa99c37b029d000001', NULL, '54814817f4e87273d2000002', '', ' ', '', '', '{}', true, false, 'Untitled-5368c1aa99c37b029d000001', false, false, false, 0, 0, 0, false, 0, '2014-12-05 05:54:51.599+00', '2014-12-05 05:54:56.678+00', NULL, NULL, '5368c1aa99c37b029d000001', 200043, true) ON CONFLICT DO NOTHING;
INSERT INTO public.notes (id, user_id, created_user_id, notebook_id, title, description, src, img_src, tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined, read_num, like_num, comment_num, is_markdown, attach_num, created_time, updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted) VALUES ('54814eb4f4e872189d000003', '5368c1aa99c37b029d000001', NULL, '54814eabf4e872189d000001', 'a', 'a ', '', '', '{}', true, false, 'a', false, false, false, 0, 0, 0, false, 0, '2014-12-05 06:20:39.856+00', '2014-12-05 06:20:39.856+00', NULL, NULL, '5368c1aa99c37b029d000001', 200042, true) ON CONFLICT DO NOTHING;
INSERT INTO public.notes (id, user_id, created_user_id, notebook_id, title, description, src, img_src, tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined, read_num, like_num, comment_num, is_markdown, attach_num, created_time, updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted) VALUES ('54817b33f4e8725881000000', '5368c1aa99c37b029d000001', NULL, '54815d48f4e87264d3000000', '', ' ', '', '', '{}', true, false, 'Untitled-5368c1aa99c37b029d000001-2', false, false, false, 0, 0, 0, false, 0, '2014-12-05 09:30:33.312+00', '2014-12-05 09:30:33.887+00', NULL, NULL, '5368c1aa99c37b029d000001', 200041, true) ON CONFLICT DO NOTHING;
INSERT INTO public.notes (id, user_id, created_user_id, notebook_id, title, description, src, img_src, tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined, read_num, like_num, comment_num, is_markdown, attach_num, created_time, updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted) VALUES ('54817b49f4e8725881000001', '5368c1aa99c37b029d000001', NULL, '548125adf4e872105c000007', '', ' ', '', '', '{}', true, false, 'Untitled-5368c1aa99c37b029d000001-3', false, false, false, 0, 0, 0, false, 0, '2014-12-05 09:30:52.615+00', '2014-12-06 15:27:35.432+00', NULL, NULL, '5368c1aa99c37b029d000001', 200038, true) ON CONFLICT DO NOTHING;
INSERT INTO public.notes (id, user_id, created_user_id, notebook_id, title, description, src, img_src, tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined, read_num, like_num, comment_num, is_markdown, attach_num, created_time, updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted) VALUES ('5482b541f4e87253cb000000', '5368c1aa99c37b029d000001', NULL, '54815d48f4e87264d3000000', 'asdf', 'asdfad ', '', '', '{}', true, false, 'asdf', false, false, false, 0, 0, 0, false, 0, '2014-12-06 07:50:28.261+00', '2014-12-06 07:50:28.261+00', NULL, NULL, '5368c1aa99c37b029d000001', 200040, true) ON CONFLICT DO NOTHING;
INSERT INTO public.notes (id, user_id, created_user_id, notebook_id, title, description, src, img_src, tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined, read_num, like_num, comment_num, is_markdown, attach_num, created_time, updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted) VALUES ('5482b54cf4e87253cb000001', '5368c1aa99c37b029d000001', NULL, '54815d48f4e87264d3000000', '', ' ', '', '', '{}', true, false, 'Untitled-5368c1aa99c37b029d000001-4', false, false, false, 0, 0, 0, false, 0, '2014-12-06 07:50:42.756+00', '2014-12-06 07:50:43.699+00', NULL, NULL, '5368c1aa99c37b029d000001', 200039, true) ON CONFLICT DO NOTHING;
INSERT INTO public.notes (id, user_id, created_user_id, notebook_id, title, description, src, img_src, tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined, read_num, like_num, comment_num, is_markdown, attach_num, created_time, updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted) VALUES ('54832064f4e87203a4000000', '5368c1aa99c37b029d000001', NULL, '548125adf4e872105c000007', '', ' ', '', '', '{}', true, false, 'Untitled-5368c1aa99c37b029d000001-5', false, false, false, 0, 0, 0, false, 0, '2014-12-06 15:27:34.166+00', '2014-12-06 15:27:56.647+00', NULL, NULL, '5368c1aa99c37b029d000001', 200037, true) ON CONFLICT DO NOTHING;
INSERT INTO public.notes (id, user_id, created_user_id, notebook_id, title, description, src, img_src, tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined, read_num, like_num, comment_num, is_markdown, attach_num, created_time, updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted) VALUES ('5483207cf4e87203a4000001', '5368c1aa99c37b029d000001', NULL, '548125adf4e872105c000007', 'About Gemsnote', '





About Gemsnote



1. Introduction

Gemsnote, not just a notepad!


Some Features


Knowledge: Manage your knowledge in gemsnote. gemsnote contains the tinymce editor and a markdown editor, just enjoy yourself writing.
Share: Share your ', '', 'http://7xj51o.com1.z0.glb.clouddn.com/default_markdown.png', '{gemsnote}', false, true, 'about-gemsnote', false, false, false, 4, 0, 0, true, 0, '2014-12-06 15:28:01.815+00', '2015-06-15 10:46:41.926+00', NULL, '2015-06-15 10:46:51.268+00', '5368c1aa99c37b029d000001', 200030, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notes (id, user_id, created_user_id, notebook_id, title, description, src, img_src, tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined, read_num, like_num, comment_num, is_markdown, attach_num, created_time, updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted) VALUES ('5524ba2f99c37b292000000b', '5524ba2f99c37b2920000007', NULL, '5524ba2f99c37b2920000008', 'Gemsnote产品说明', 'Gemsnote产品说明

1. 介绍

Gemsnote, 不只是笔记!

特性


知识管理: 通过gemsnote来管理知识, gemsnote有易操作的界面, 包含两款编辑器tinymce和markdown. 在gemsnote, 你可以尽情享受写作.
分享: 你也可以通过分享知识给好友, 让好友拥有你的知识.
协作: 在分享的同时也可以与好友一起协作知识.
博客: gemsnote也可以作为你的博客, 将知识公开成博客, 让gemsnote把你的知识传播的更远!


2. 为什么我们要创建gemsnote?

说实话, 我们曾是evernote的忠实粉丝, 但是我们也发现evernote的不足:', '', '', '{}', false, false, 'Gemsnote%E4%BA%A7%E5%93%81%E8%AF%B4%E6%98%8E', false, false, false, 0, 0, 0, true, 0, '2015-04-08 05:18:39.431+00', '2015-04-08 05:18:39.437+00', NULL, '2015-04-08 05:18:39.431+00', '5524ba2f99c37b2920000007', 5, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notes (id, user_id, created_user_id, notebook_id, title, description, src, img_src, tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined, read_num, like_num, comment_num, is_markdown, attach_num, created_time, updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted) VALUES ('557eaa9705fcd14d95000000', '5368c1aa99c37b029d000001', NULL, '548125adf4e872105c000007', '', ' ', '', '', '{""}', true, false, 'Untitled-5368c1aa99c37b029d000001-7', false, false, false, 0, 0, 0, false, 0, '2015-06-15 10:36:09.572+00', '2015-06-15 10:36:09.572+00', NULL, '2015-06-15 10:36:09.572+00', '5368c1aa99c37b029d000001', 200036, true) ON CONFLICT DO NOTHING;
INSERT INTO public.notes (id, user_id, created_user_id, notebook_id, title, description, src, img_src, tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined, read_num, like_num, comment_num, is_markdown, attach_num, created_time, updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted) VALUES ('557eaa9905fcd14d95000001', '5368c1aa99c37b029d000001', NULL, '548125adf4e872105c000007', '关于Gemsnote', '





关于Gemsnote



1. 介绍

Gemsnote, 不只是笔记!


特性


知识管理: 通过gemsnote来管理知识, gemsnote有易操作的界面, 包含两款编辑器tinymce和markdown. 在gemsnote, 你可以尽情享受写作.
分享: 你也可以通过分享知识给好友, 让好友拥有你的知识.
协作: 在分享的同时也可以与好友一起协作知识.
博客: gemsnote也可以作为你的博客, 将知识公开成博客, 让gemsnote把你的知识传播的更远!', '', 'http://7xj51o.com1.z0.glb.clouddn.com/default_markdown.png', '{""}', false, false, 'Untitled-5368c1aa99c37b029d000001-8', false, false, false, 0, 0, 0, true, 0, '2015-06-15 10:36:13.324+00', '2015-06-15 10:37:24.144+00', NULL, '2015-06-15 10:36:13.324+00', '5368c1aa99c37b029d000001', 200014, false) ON CONFLICT DO NOTHING;
INSERT INTO public.notes (id, user_id, created_user_id, notebook_id, title, description, src, img_src, tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined, read_num, like_num, comment_num, is_markdown, attach_num, created_time, updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted) VALUES ('561bb57505fcd164d3000000', '5368c1aa99c37b029d000001', NULL, '557eab5705fcd14d95000002', 'Hello', 'dddd ', '', '', '{""}', true, false, 'Hello', false, false, false, 0, 0, 0, false, 0, '2015-10-12 13:28:25.575+00', '2015-10-12 13:28:26.955+00', NULL, '2015-10-12 13:28:25.575+00', '5368c1aa99c37b029d000001', 200034, true) ON CONFLICT DO NOTHING;


--
-- Data for Name: gemsnote_schema_migrations; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: share_notebooks; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.share_notebooks (id, user_id, to_user_id, to_group_id, to_group, notebook_id, seq, perm, created_time) VALUES ('5368c9fcdd73244c2b52e10b', '5368c1aa99c37b029d000001', '5368c9fc99c37b095a000006', NULL, NULL, '5368c1aa99c37b029d000002', 0, 0, '2014-05-06 11:39:40.634+00') ON CONFLICT DO NOTHING;
INSERT INTO public.share_notebooks (id, user_id, to_user_id, to_group_id, to_group, notebook_id, seq, perm, created_time) VALUES ('540814f1acf7541eaf8369c8', '5368c1aa99c37b029d000001', '540814f199c37b555d000001', NULL, NULL, '5368c1aa99c37b029d000002', 0, 0, '2014-09-04 07:29:53.184+00') ON CONFLICT DO NOTHING;
INSERT INTO public.share_notebooks (id, user_id, to_user_id, to_group_id, to_group, notebook_id, seq, perm, created_time) VALUES ('5463264d8a88dcbd99de2796', '5368c1aa99c37b029d000001', NULL, '5463263299c37b80ae000009', NULL, '54479ae219807a4cef000000', 0, 1, '2014-11-12 09:20:13.267+00') ON CONFLICT DO NOTHING;
INSERT INTO public.share_notebooks (id, user_id, to_user_id, to_group_id, to_group, notebook_id, seq, perm, created_time) VALUES ('5524b99b82d7216d0e5516c1', '5368c1aa99c37b029d000001', '5524b99b99c37b2920000002', NULL, NULL, '54479ae219807a4cef000000', 0, 1, '2015-04-08 05:16:11.709+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: share_notes; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.share_notes (id, user_id, to_user_id, to_group_id, to_group, note_id, perm, created_time) VALUES ('5524b99b82d7216d0e5516c3', '5368c1aa99c37b029d000001', '5524b99b99c37b2920000002', NULL, NULL, '5368c1b919807a6f95000000', 1, '2015-04-08 05:16:11.724+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: tag_count; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.tag_count (id, user_id, tag, is_blog, count) VALUES ('557ead1ba055f4118f195db7', '5368c1aa99c37b029d000001', 'gemsnote', true, 1) ON CONFLICT DO NOTHING;


--
-- Data for Name: tags; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.tags (id, tags) VALUES ('5368c1aa99c37b029d000001', '{gemsnote,欢迎,red}') ON CONFLICT DO NOTHING;


--
-- Data for Name: themes; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.themes (id, user_id, name, version, author, author_url, path, info, is_active, is_default, style, created_time, updated_time) VALUES ('54633b3d99c37b969d000001', '5368c1aa99c37b029d000001', 'gemsnote default theme', '1.0', 'gemsnote.com', 'http://gemsnote.com', 'public/blog/themes/default', '{"Name": "gemsnote default theme", "Author": "gemsnote.com", "Version": "1.0", "AuthorUrl": "http://gemsnote.com", "FriendLinks": [{"Url": "http://gemsnote.com/note", "Title": "我的笔记"}, {"Url": "http://gemsnote.com", "Title": "gemsnote home"}, {"Url": "http://bbs.gemsnote.com", "Title": "gemsnote 社区"}, {"Url": "http://lea.gemsnote.com", "Title": "lea++"}, {"Url": "https://github.com/gemsnote/gemsnote", "Title": "gemsnote github"}]}', false, true, 'blog_default', '2014-11-12 10:49:33.878+00', '2014-11-12 10:49:33.878+00') ON CONFLICT DO NOTHING;
INSERT INTO public.themes (id, user_id, name, version, author, author_url, path, info, is_active, is_default, style, created_time, updated_time) VALUES ('54633b3d99c37b969d000002', '5368c1aa99c37b029d000001', 'gemsnote elegant', '1.0', 'gemsnote.com', 'http://gemsnote.com', 'public/blog/themes/elegant', '{"Name": "gemsnote elegant", "Author": "gemsnote.com", "Version": "1.0", "AuthorUrl": "http://gemsnote.com", "FriendLinks": [{"Url": "http://gemsnote.com/note", "Title": "我的笔记"}, {"Url": "http://gemsnote.com", "Title": "gemsnote home"}, {"Url": "http://bbs.gemsnote.com", "Title": "gemsnote 社区"}, {"Url": "http://lea.gemsnote.com", "Title": "lea++"}, {"Url": "https://github.com/gemsnote/gemsnote", "Title": "gemsnote github"}]}', true, true, 'blog_daqi', '2014-11-12 10:49:33.882+00', '2014-11-12 10:49:33.882+00') ON CONFLICT DO NOTHING;
INSERT INTO public.themes (id, user_id, name, version, author, author_url, path, info, is_active, is_default, style, created_time, updated_time) VALUES ('54633b3d99c37b969d000003', '5368c1aa99c37b029d000001', 'gemsnote nav fixed', '1.0', 'gemsnote.com', 'http://gemsnote.com', 'public/blog/themes/nav_fixed', '{"Name": "gemsnote nav fixed", "Author": "gemsnote.com", "Version": "1.0", "AuthorUrl": "http://gemsnote.com", "FriendLinks": [{"Url": "http://gemsnote.com/note", "Title": "我的笔记"}, {"Url": "http://gemsnote.com", "Title": "gemsnote home"}, {"Url": "http://bbs.gemsnote.com", "Title": "gemsnote 社区"}, {"Url": "http://lea.gemsnote.com", "Title": "lea++"}, {"Url": "https://github.com/gemsnote/gemsnote", "Title": "gemsnote github"}]}', false, true, 'blog_left_fixed', '2014-11-12 10:49:33.882+00', '2014-11-12 10:49:33.882+00') ON CONFLICT DO NOTHING;
INSERT INTO public.themes (id, user_id, name, version, author, author_url, path, info, is_active, is_default, style, created_time, updated_time) VALUES ('5524ba4499c37b292000000d', '5524ba2f99c37b2920000007', 'Gemsnote default theme', '1.0', 'gemsnote.com', 'http://gemsnote.com', 'public/upload/5524ba2f99c37b2920000007/themes/5524ba4499c37b292000000d', '{"Desc": "", "Name": "Gemsnote default theme", "Author": "gemsnote.com", "Version": "1.0", "AuthorUrl": "http://gemsnote.com", "FriendLinks": [{"Url": "http://gemsnote.com/note", "Title": "My Note"}, {"Url": "http://gemsnote.com", "Title": "Gemsnote Home"}, {"Url": "https://groups.google.com/forum/?fromgroups#!forum/gemsnote", "Title": "Gemsnote Comunity"}, {"Url": "http://lea.gemsnote.com", "Title": "lea++"}, {"Url": "https://github.com/gemsnote/gemsnote", "Title": "Gemsnote Github"}]}', false, false, '', '2015-04-08 05:19:00.553+00', '2015-04-08 05:19:00.553+00') ON CONFLICT DO NOTHING;
INSERT INTO public.themes (id, user_id, name, version, author, author_url, path, info, is_active, is_default, style, created_time, updated_time) VALUES ('5524ba4499c37b292000000e', '5524ba2f99c37b2920000007', 'Gemsnote elegant', '1.0', 'gemsnote.com', 'http://gemsnote.com', 'public/upload/5524ba2f99c37b2920000007/themes/5524ba4499c37b292000000e', '{"Desc": "", "Name": "Gemsnote elegant", "Author": "gemsnote.com", "Version": "1.0", "AuthorUrl": "http://gemsnote.com", "FriendLinks": [{"Url": "http://gemsnote.com/note", "Title": "My Note"}, {"Url": "http://gemsnote.com", "Title": "Gemsnote Home"}, {"Url": "https://groups.google.com/forum/?fromgroups#!forum/gemsnote", "Title": "Gemsnote Comunity"}, {"Url": "http://lea.gemsnote.com", "Title": "lea++"}, {"Url": "https://github.com/gemsnote/gemsnote", "Title": "Gemsnote Github"}]}', true, false, '', '2015-04-08 05:19:00.579+00', '2015-04-08 05:19:00.579+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: user_blogs; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.user_blogs (id, logo, title, sub_title, about_me, can_comment, comment_type, disqus_id, style, css, theme_id, theme_path, cate_ids, singles, per_page_size, sort_field, is_asc, sub_domain, domain) VALUES ('5368c1aa99c37b029d000001', '', 'Gemsnote''s Blog', 'I love Gemsnote!', '<p>Hello, 大家好, 我是gemsnote, 赶紧来体验gemsnote吧!!</p>', true, 'default', 'gemsnote', 'blog_daqi', '', '54633b3d99c37b969d000002', '', '{}', '[{"Title": "About Me", "SingleId": "546325e699c37b80ae000001", "UrlTitle": "About-Me"}]', 0, '', false, '', '') ON CONFLICT DO NOTHING;
INSERT INTO public.user_blogs (id, logo, title, sub_title, about_me, can_comment, comment_type, disqus_id, style, css, theme_id, theme_path, cate_ids, singles, per_page_size, sort_field, is_asc, sub_domain, domain) VALUES ('5368c9fc99c37b095a000006', '', '我的博客', 'love gemsnote!', '<p>Hello, I am (^_^)</p>', false, '', '', 'blog_left_fixed', '', NULL, '', '{}', '[{"Title": "About Me", "SingleId": "546325e699c37b80ae000002", "UrlTitle": "About-Me"}]', 0, '', false, '', '') ON CONFLICT DO NOTHING;
INSERT INTO public.user_blogs (id, logo, title, sub_title, about_me, can_comment, comment_type, disqus_id, style, css, theme_id, theme_path, cate_ids, singles, per_page_size, sort_field, is_asc, sub_domain, domain) VALUES ('5524ba2f99c37b2920000007', '', 'b@a.com ''s Blog', 'Love Gemsnote!', 'Hello, I am (^_^)', true, '', '', '', '', '5524ba4499c37b292000000e', '', '{}', '[{"Title": "About Me", "SingleId": "5524ba2f99c37b292000000c", "UrlTitle": "About-Me"}]', 0, '', false, '', '') ON CONFLICT DO NOTHING;


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.users (id, email, verified, username, username_raw, pwd, created_time, logo, theme, notebook_width, note_list_width, md_editor_width, left_is_min, third_user_id, third_username, third_type, image_num, image_size, attach_num, attach_size, from_user_id, account_type, account_start_time, account_end_time, max_image_nums, max_image_size, max_attach_num, max_attach_size, max_per_attach_size, usn, full_sync_before) VALUES ('5368c1aa99c37b029d000001', 'admin@gemsnote.com', false, 'admin', 'admin', '4533cd5ac65d7b5ccb02f396ff4b6615', '2014-05-06 11:04:10.658+00', '', 'simple', 160, 266, 0, false, '', '', 0, 0, 0, 0, 0, NULL, '', NULL, NULL, 0, 0, 0, 0, 0, 200044, NULL) ON CONFLICT DO NOTHING;


--
-- PostgreSQL database dump complete
--

-- Gemsnote welcome notes for new installations.
UPDATE public.notes SET title = 'About Gemsnote', description = 'Gemsnote features and deployment overview', src = '', img_src = '', tags = '{gemsnote}', is_blog = false, url_title = 'about-gemsnote', is_markdown = true WHERE id = '5483207cf4e87203a4000001';
UPDATE public.note_contents SET is_blog = false, content = $gemsnote_en$# About Gemsnote

Gemsnote originates from Gemsnote and is a self-hosted knowledge and notes service.

## Features

- Choose MongoDB or PostgreSQL through configuration.
- Connect directly to an existing Gemsnote MongoDB database, or migrate it to PostgreSQL.
- Use the Web interface or the Desktop client.
- Create rich-text and Markdown notes.
- Synchronize across devices, share notes, and cache shared notes for offline Desktop use.
- Manage accounts and the server from the administration interface.
- Blog functionality is intentionally not included.

Gemsnote is designed to run on your own server, so your data remains under your control.$gemsnote_en$, abstract = 'Gemsnote features and deployment overview' WHERE id = '5483207cf4e87203a4000001';
UPDATE public.note_content_histories SET histories = jsonb_build_array(jsonb_build_object('Content', (SELECT content FROM public.note_contents WHERE id = '5483207cf4e87203a4000001'), 'UpdatedTime', '2015-06-15T18:37:30.724+08:00', 'UpdatedUserId', user_id)) WHERE id = '5483207cf4e87203a4000001';

UPDATE public.notes SET title = '关于 Gemsnote', description = 'Gemsnote 功能与部署简介', src = '', img_src = '', tags = '{gemsnote}', is_blog = false, url_title = 'about-gemsnote-zh', is_markdown = true WHERE id = '557eaa9905fcd14d95000001';
UPDATE public.note_contents SET is_blog = false, content = $gemsnote_zh$# 关于 Gemsnote

Gemsnote 源自 Gemsnote，是一个可自托管的知识与笔记服务。

## 功能

- 可通过配置选择 MongoDB 或 PostgreSQL。
- 可直接连接已有的 Gemsnote MongoDB，或将数据迁移到 PostgreSQL。
- 提供 Web 界面和 Desktop 客户端。
- 支持富文本与 Markdown 笔记。
- 支持多端同步、笔记共享，以及 Desktop 端共享笔记离线缓存。
- 提供账号管理和系统管理。
- 不再提供博客功能。

Gemsnote 适合部署在自己的服务器上，让数据始终由你掌控。$gemsnote_zh$, abstract = 'Gemsnote 功能与部署简介' WHERE id = '557eaa9905fcd14d95000001';
UPDATE public.note_content_histories SET histories = jsonb_build_array(jsonb_build_object('Content', (SELECT content FROM public.note_contents WHERE id = '557eaa9905fcd14d95000001'), 'UpdatedTime', '2015-06-15T18:37:24.15+08:00', 'UpdatedUserId', user_id)) WHERE id = '557eaa9905fcd14d95000001';

COMMIT;
