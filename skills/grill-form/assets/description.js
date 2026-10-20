// Rich descriptions from jstack 7a7c919; original Mermaid Tiny renderer.
export function element(tag,text,className){const node=document.createElement(tag);if(text!==undefined)node.textContent=text;if(className)node.className=className;return node}
const escapeHTML=text=>text.replace(/[&<>"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));
let diagramNumber=0;
let diagramQueue=Promise.resolve();
if(globalThis.mermaid)mermaid.initialize({startOnLoad:false,securityLevel:'strict',theme:'dark',fontFamily:'Arial, sans-serif',htmlLabels:false,flowchart:{htmlLabels:false},suppressErrorRendering:true,maxTextSize:12000,maxEdges:150});
function diagramFallback(container,source,reason){
 container.replaceChildren(element('p','Diagram could not be rendered. '+reason,'diagram-error'));
 const details=element('details');details.open=true;
 const pre=element('pre');pre.append(element('code',source));
 details.append(element('summary','Diagram source'),pre);container.append(details);container.dataset.renderState='error';
}
async function renderDiagram(container,source){
 const first=source.split(/\r?\n/).map(line=>line.trim()).find(line=>line&&!line.startsWith('%%'))||'';
 if(/%%\s*\{/.test(source)||! /^(?:(?:flowchart|graph)\s+(?:TB|TD|BT|RL|LR)\b|sequenceDiagram\b|stateDiagram(?:-v2)?\b)/.test(first)){
  diagramFallback(container,source,'Use a flowchart, sequence diagram, or state diagram without configuration directives.');return;
 }
 try{
  const {svg}=await mermaid.render('grill-diagram-'+(++diagramNumber),source);
  container.innerHTML=DOMPurify.sanitize(svg,{USE_PROFILES:{svg:true,svgFilters:true},FORBID_TAGS:['foreignObject','image','a']});
  const drawing=container.querySelector('svg');
  const width=drawing?.viewBox.baseVal.width;
  if(width){drawing.style.width=Math.ceil(width)+'px';drawing.style.maxWidth='';}
  container.dataset.renderState='ready';
 }catch{diagramFallback(container,source,'You can still answer this question.');}
}
export function description(source,number){
 const wrapper=element('section',undefined,'description');wrapper.setAttribute('aria-label','Description for question '+number);
 const toolbar=element('div',undefined,'description-toolbar');
 const expand=element('button',undefined,'expand');expand.type='button';expand.setAttribute('aria-expanded','false');
 const icon=document.createElementNS('http://www.w3.org/2000/svg','svg');icon.setAttribute('viewBox','0 0 24 24');icon.setAttribute('fill','none');icon.setAttribute('stroke','currentColor');icon.setAttribute('stroke-width','1.75');icon.setAttribute('stroke-linecap','round');icon.setAttribute('stroke-linejoin','round');icon.setAttribute('aria-hidden','true');
 const iconPath=document.createElementNS('http://www.w3.org/2000/svg','path');icon.append(iconPath);expand.append(icon);
 function setExpanded(expanded){const label=expanded?'Collapse description':'Expand description';expand.setAttribute('aria-expanded',String(expanded));expand.setAttribute('aria-label',label);expand.title=label;iconPath.setAttribute('d',expanded?'M3 9h6V3M9 9 3 3M21 15h-6v6M15 15l6 6':'M9 3H3v6M3 3l6 6M15 21h6v-6M21 21l-6-6');}
 setExpanded(false);
 const content=element('div',undefined,'rich-content');content.id='description-'+number;content.tabIndex=0;content.setAttribute('role','region');content.setAttribute('aria-label','Question '+number+' description');expand.setAttribute('aria-controls',content.id);
 expand.addEventListener('click',()=>setExpanded(wrapper.classList.toggle('is-expanded')));
 wrapper.addEventListener('keydown',event=>{if(event.key==='Escape'&&wrapper.classList.contains('is-expanded')){expand.click();expand.focus();}});
 toolbar.append(expand);wrapper.append(content,toolbar);
 try{
  const renderer=new marked.Renderer();
  renderer.html=({text})=>escapeHTML(text);
  renderer.image=({text})=>escapeHTML(text);
  const html=marked.parse(source,{gfm:true,breaks:false,renderer});
  content.innerHTML=DOMPurify.sanitize(html,{ALLOWED_TAGS:['p','br','strong','em','del','a','ul','ol','li','blockquote','h1','h2','h3','h4','h5','h6','pre','code','table','thead','tbody','tr','th','td','hr'],ALLOWED_ATTR:['href','title','class','align','start']});
  content.querySelectorAll('a').forEach(link=>{if(!/^https?:\/\//i.test(link.getAttribute('href')||'')){link.removeAttribute('href');return;}link.target='_blank';link.rel='noopener noreferrer';});
  content.querySelectorAll('pre > code.language-mermaid').forEach(code=>{
   const diagram=element('div','Rendering diagram…','diagram');diagram.dataset.renderState='pending';diagram.tabIndex=0;diagram.setAttribute('role','region');diagram.setAttribute('aria-label','Diagram for question '+number);
   const definition=code.textContent;code.parentElement.replaceWith(diagram);
   diagramQueue=diagramQueue.then(()=>renderDiagram(diagram,definition)).catch(()=>diagramFallback(diagram,definition,'You can still answer this question.'));
  });
 }catch{content.append(element('p','Description could not be rendered; original source follows.','error'),element('pre',source));}
 return wrapper;
}
