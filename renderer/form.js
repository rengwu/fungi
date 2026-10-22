(async()=>{
  const ready=await fungi.ready();
  const schema=ready.form;
  document.title=schema.title;
  document.querySelector('#title').textContent=schema.title;
  document.querySelector('#context').textContent=schema.context||'';
  if(ready.delivery?.mode==='manual')document.querySelector('#delivery').textContent='After sending, tell your agent “done” so it can collect your response.';
  const form=document.querySelector('#form'), fields=document.querySelector('#fields'), error=document.querySelector('#error'), status=document.querySelector('#status');
  fields.style.display='grid';fields.style.gap='24px';
  const values=Object.create(null), images=new Map();
  let saving=Promise.resolve(), generation=0, busy=false;
  function imageList(){return [...images].map(([key,value])=>({key,...value}));}
  function save(){
    const n=++generation, data=structuredClone(values), files=imageList();
    status.textContent='Saving draft…';
    saving=saving.catch(()=>{}).then(()=>fungi.saveDraft(data,files)).then(()=>{
      if(n===generation)status.textContent='Draft saved';
    }).catch(e=>{error.textContent=e.message;status.textContent='Draft not saved';throw e;});
    // This handler observes rejection; submit still awaits the original promise.
    saving.catch(()=>{});
  }
  for(const field of schema.fields){
    const box=document.createElement('div'),label=document.createElement('label');
    label.textContent=field.label||field.name;label.htmlFor='field-'+field.name;box.append(label);
    values[field.name]=ready.draft?.[field.name] ?? ready.input?.[field.name] ?? field.value ?? (field.type==='range'?field.min:field.type==='image'?null:'');
    if(field.type==='image'){
      box.classList.add('image');box.tabIndex=0;
      const help=document.createElement('p');help.textContent='Choose an image, or focus this area and paste one.';
      const input=document.createElement('input');input.type='file';input.accept='image/*';input.id=label.htmlFor;
      const pasteArea=document.createElement('div');pasteArea.contentEditable='true';pasteArea.role='textbox';
      pasteArea.setAttribute('aria-label','Paste image for '+(field.label||field.name));pasteArea.className='paste-area';
      pasteArea.textContent='Click here and paste an image';
      pasteArea.oninput=()=>{pasteArea.textContent='Click here and paste an image';};
      const preview=document.createElement('img');preview.hidden=true;preview.alt='Selected reference';
      const remove=document.createElement('button');remove.type='button';remove.className='remove';remove.textContent='Remove image';remove.hidden=true;
      let url;
      function display(blob){
        if(url)URL.revokeObjectURL(url);
        preview.hidden=!blob;remove.hidden=!blob;
        if(blob){url=URL.createObjectURL(blob);preview.src=url;}
      }
      async function select(blob,name){
        if(!blob.type.startsWith('image/')){error.textContent='Please choose an image.';return;}
        if(blob.size>ready.limits.artifactBytes){error.textContent='Image exceeds 20 MiB.';return;}
        error.textContent='';
        images.set(field.name,{name,blob});values[field.name]={artifact:field.name};display(blob);save();
      }
      const restored=ready.attachments.find(a=>a.key===field.name);
      if(restored){images.set(field.name,{name:restored.name,blob:restored.blob});values[field.name]={artifact:field.name};display(restored.blob);}
      else values[field.name]=null;
      input.onchange=()=>{if(input.files[0])select(input.files[0],input.files[0].name);};
      box.onpaste=e=>{
        const file=[...e.clipboardData.items].find(i=>i.kind==='file'&&i.type.startsWith('image/'))?.getAsFile();
        if(file){e.preventDefault();select(file,file.name||'pasted-image.png');}
        else if(e.target===pasteArea){e.preventDefault();error.textContent='The clipboard does not contain an image.';}
      };
      remove.onclick=()=>{images.delete(field.name);values[field.name]=null;input.value='';display(null);save();};
      box.append(help,input,pasteArea,preview,remove);
    }else{
      const input=document.createElement(field.type==='textarea'?'textarea':'input');
      input.id=label.htmlFor;input.name=field.name;input.required=!!field.required;
      if(field.type==='range'){
        input.type='range';input.min=field.min;input.max=field.max;input.step=field.step??1;
        const output=document.createElement('output');label.append(output);
        input.value=values[field.name];values[field.name]=Number(input.value);output.value=input.value;
        input.oninput=()=>{values[field.name]=Number(input.value);output.value=input.value;save();};
      }else{
        input.type='text';input.value=values[field.name];
        input.oninput=()=>{values[field.name]=input.value;save();};
      }
      box.append(input);
    }
    fields.append(box);
  }
  function setBusy(value){busy=value;form.inert=value;for(const button of form.querySelectorAll('button'))button.disabled=value;}
  form.onsubmit=async event=>{
    event.preventDefault();if(busy)return;error.textContent='';setBusy(true);
    try{
      for(const field of schema.fields)if(field.type==='image'&&field.required&&!images.has(field.name))throw new Error((field.label||field.name)+' is required');
      // The final submission retries current data even if an earlier autosave failed.
      await saving.catch(()=>{});
      status.textContent='Sending…';
      await fungi.submit({...values},imageList());
    }catch(e){error.textContent=e.message;status.textContent='Response not sent';setBusy(false);}
  };
  document.querySelector('#cancel').onclick=async()=>{
    if(busy)return;setBusy(true);
    try{await saving.catch(()=>{});await fungi.cancel();}
    catch(e){error.textContent=e.message;setBusy(false);}
  };
  // Save resolved defaults/input too, so a close followed by reopen is stable.
  save();
})().catch(error=>{document.querySelector('#error').textContent=error.message;});
