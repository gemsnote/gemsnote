// Rejection sampling avoids modulo bias. Never fall back to Math.random.
function randomIndex(size: number): number {
  const value = new Uint8Array(1)
  const limit = 256 - (256 % size)
  do { crypto.getRandomValues(value) } while (value[0] >= limit)
  return value[0] % size
}

export function generateInitialPassword(): string {
  const groups = ['ABCDEFGHJKLMNPQRSTUVWXYZ', 'abcdefghijkmnopqrstuvwxyz', '23456789', '!@#$%^&*_-+=?']
  const alphabet = groups.join('')
  const chars = groups.map(group => group[randomIndex(group.length)])
  while (chars.length < 12) chars.push(alphabet[randomIndex(alphabet.length)])
  for (let i = chars.length - 1; i > 0; i--) {
    const j = randomIndex(i + 1)
    ;[chars[i], chars[j]] = [chars[j], chars[i]]
  }
  return chars.join('')
}
