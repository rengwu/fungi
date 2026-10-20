import {VERSION, normalizeSpec, answered, validateAnswers, restoreAnswers} from './model.mjs';
import {element, description} from './description.js';

const form=document.querySelector('#questions'), reference=document.querySelector('#reference');
const submit=document.querySelector('#submit'), cancel=document.querySelector('#cancel');
const status=document.querySelector('#status'), message=document.querySelector('#message');
function showError(error){message.textContent=error.message;message.hidden=false;}

(async()=>{
  const ready=await fungi.ready();
  const spec=normalizeSpec(ready.input,{richText:true});
  const projectName=typeof ready.input.projectName==='string' ? ready.input.projectName : 'Fungi';
  const answers=restoreAnswers(spec,ready.draft);
  document.querySelector('#project-name').textContent=projectName;
  document.querySelector('#title').textContent=spec.title;
  document.title=projectName+' — '+spec.title;
  if(ready.delivery?.mode==='manual')document.querySelector('#delivery').textContent='After submitting, tell your agent “done” so it can collect your answers.';

  let image=ready.attachments.find(a=>a.key==='reference'), previewURL;
  let saving=Promise.resolve(), generation=0, busy=false;
  const files=()=>image ? [image] : [];
  function updateSubmit(){submit.disabled=busy||!answers.every(answered);}
  function save(){
    updateSubmit();message.hidden=true;
    const n=++generation;
    const draft={schemaVersion:VERSION,spec,answers:structuredClone(answers)};
    const attachments=files();
    status.textContent='Saving draft…';
    saving=saving.catch(()=>{}).then(()=>fungi.saveDraft(draft,attachments)).then(()=>{
      if(n===generation)status.textContent='Draft saved';
    }).catch(error=>{
      if(n===generation){status.textContent='Draft not saved';showError(error);}
      throw error;
    });
    saving.catch(()=>{}); // Keep failures observable to the UI without an unhandled rejection.
  }
  for(const [i,q] of spec.questions.entries()){
    const fieldset=element('fieldset');
    fieldset.append(element('legend',q.number+'. '+q.prompt));
    if(q.context.trim())fieldset.append(description(q.context,q.number));
    fieldset.append(element('p','Recommended: '+q.recommendation,'recommendation'));
    const choices=element('div',undefined,'choices');
    for(const option of [...q.options,{id:'other',label:'Something else'}]){
      const label=element('label',undefined,'choice'), radio=element('input');
      radio.type='radio';radio.name=q.id;radio.value=option.id;radio.checked=answers[i].choiceId===option.id;
      radio.onchange=()=>{answers[i].choiceId=option.id;save();};
      label.append(radio,element('span',option.label));choices.append(label);
    }
    const label=element('label','Notes','notes-label'), notes=element('textarea');
    label.htmlFor=notes.id='notes-'+q.id;notes.rows=3;notes.maxLength=12000;
    notes.setAttribute('aria-label','Notes for question '+q.number);notes.value=answers[i].notes;
    notes.oninput=()=>{answers[i].notes=notes.value;save();};
    fieldset.append(choices,label,notes);form.insertBefore(fieldset,reference);
  }

  const picker=document.querySelector('#image'), preview=document.querySelector('#preview');
  const remove=document.querySelector('#remove-image'), filename=document.querySelector('#attachment-name');
  function displayImage(){
    if(previewURL)URL.revokeObjectURL(previewURL);
    preview.hidden=remove.hidden=!image;filename.textContent=image?.name||'';
    if(image){previewURL=URL.createObjectURL(image.blob);preview.src=previewURL;}
    else preview.removeAttribute('src');
  }
  picker.onchange=()=>{
    const blob=picker.files[0];if(!blob)return;
    if(!blob.type.startsWith('image/')){showError(new Error('Please choose an image.'));picker.value='';return;}
    if(blob.size>ready.limits.artifactBytes){showError(new Error('Image exceeds 20 MiB.'));picker.value='';return;}
    image={key:'reference',name:blob.name,blob};displayImage();save();
  };
  remove.onclick=()=>{image=undefined;picker.value='';displayImage();save();};
  reference.hidden=false;displayImage();cancel.disabled=false;

  function setBusy(value){
    busy=value;
    for(const control of form.querySelectorAll('input,textarea,button'))control.disabled=value;
    updateSubmit();
  }
  form.onsubmit=async event=>{
    event.preventDefault();if(busy||!answers.every(answered))return;
    const data={schemaVersion:VERSION,projectName,title:spec.title,answers:validateAnswers(spec,answers),
      reference:image?{artifact:'reference'}:null};
    setBusy(true);message.hidden=true;
    try{
      await saving.catch(()=>{});
      status.textContent='Submitting…';
      await fungi.submit(data,files());
    }catch(error){status.textContent='Answers not submitted';showError(error);setBusy(false);}
  };
  cancel.onclick=async()=>{
    if(busy)return;setBusy(true);message.hidden=true;
    try{await saving.catch(()=>{});await fungi.cancel();}
    catch(error){showError(error);setBusy(false);}
  };
  save();
})().catch(error=>{status.textContent='Questionnaire could not be opened';showError(error);});
