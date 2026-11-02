import {build} from 'esbuild'
import {readFile, writeFile, mkdir, access} from 'node:fs/promises'
import {fileURLToPath} from 'node:url'
import {createHash} from 'node:crypto'
const root = new URL('../', import.meta.url)
const result = await build({
  entryPoints: [fileURLToPath(new URL('mermaid/renderer.ts', root))],
  bundle: true, write: false, minify: true, format: 'iife', target: 'es2020',
  legalComments: 'inline',
})
const pkg = JSON.parse(await readFile(new URL('node_modules/mermaid/package.json', root), 'utf8'))
// Inline the complete bundle, including lazy diagram modules, for file:// WebViews.
const js = result.outputFiles[0].text.replace(/<\/script/gi, '<\\/script')
const html = `<!doctype html><html><head><meta charset="utf-8">
<meta http-equiv="Content-Security-Policy" content="default-src 'none'; script-src 'unsafe-inline'; style-src 'unsafe-inline'; img-src data:; font-src 'none'; connect-src 'none'; base-uri 'none'; form-action 'none'">
<style>body{margin:0;background:white}#graph{width:1024px}svg{overflow:visible}</style>
</head><body><div id="graph"></div><script>${js}</script></body></html>`
const license = await readFile(new URL('node_modules/mermaid/LICENSE', root), 'utf8')
const manifest = JSON.stringify({mermaid: pkg.version, sha256: createHash('sha256').update(html).digest('hex')}, null, 2) + '\n'
const targets = [new URL('public/assets/mermaid/', root)]
if (process.argv.includes('--mobile')) {
  await access(new URL('../mobile-app/pubspec.yaml', root))
  targets.push(new URL('../mobile-app/assets/mermaid/', root))
}
for (const dir of targets) {
  await mkdir(dir, {recursive: true})
  await writeFile(new URL('index.html', dir), html)
  await writeFile(new URL('LICENSE', dir), license)
  await writeFile(new URL('manifest.json', dir), manifest)
}
console.log(`Mermaid ${pkg.version}: offline renderer built for ${targets.length} target(s)`)
