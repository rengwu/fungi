// Adapt an already generated jstack form. Keep upstream rendering untouched;
// replace browser-download delivery with Fungi's acknowledged bridge.
import {readFile, writeFile} from 'node:fs/promises';
import {fileURLToPath} from 'node:url';

const [input, output] = process.argv.slice(2);
if (!input || !output) throw new Error('Usage: node adapt.mjs GENERATED_FORM OUTPUT_HTML');
let html = await readFile(input, 'utf8');
const helpers = [];
html = html.replace(/<script>\n([\s\S]*?)<\/script>/g, (tag, code) => {
  if (code.startsWith('function validRegion(') ||
      code.startsWith('// Inline image editing') ||
      code.startsWith('// Every preview')) {
    helpers.push(code);
    return '';
  }
  return tag;
});
if (helpers.length !== 3) throw new Error('Unexpected jstack helper layout; review this adapter');
const marker = "<script>\n'use strict';";
const start = html.indexOf(marker), end = html.indexOf('</script>', start);
if (start < 0 || end < 0) throw new Error('Missing jstack form script');
let code = html.slice(start + '<script>\n'.length, end);
function replaceOnce(before, after) {
  if (code.split(before).length !== 2) throw new Error('Unexpected jstack form layout: ' + before.slice(0, 60));
  code = code.replace(before, () => after);
}
replaceOnce("JSON.parse(sessionStorage.getItem(storageKey))", 'await restoreFungiDraft(state)');
const updateStart = code.indexOf('function update(){');
const updateEnd = code.indexOf('questions.forEach((q,i)=>{', updateStart);
if (updateStart < 0 || updateEnd < 0) throw new Error('Missing update boundary');
code = code.slice(0, updateStart) + 'function update(){ saveFungiDraft(); }\n' + code.slice(updateEnd);
const submitStart = code.indexOf("form.addEventListener('submit',event=>{");
if (submitStart < 0 || !code.trimEnd().endsWith('});')) throw new Error('Unexpected submit handler');
code = code.slice(0, submitStart) + "form.addEventListener('submit',submitFungiForm);\n";
const bridge = await readFile(fileURLToPath(new URL('./bridge.js', import.meta.url)), 'utf8');
const replacement = '<script>\n(async()=>{\nconst fungiReady=await fungi.ready();\n' +
  helpers.join('\n') + '\n' + bridge + '\n' + code +
  '\ndocument.querySelector(".review-hint").textContent="Mock questionnaire · Nothing here is a project decision. Try any items, then submit to return your feedback directly to the agent.";\n' +
  '})().catch(error=>{const node=document.getElementById("message");node.textContent=error.message;node.hidden=false;});\n</script>';
html = html.slice(0, start) + replacement + html.slice(end + '</script>'.length);
await writeFile(output, html, {mode:0o600, flag:'wx'});
console.log(output);
