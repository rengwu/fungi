const fs=require('node:fs/promises');
const path=require('node:path');
const {Readable}=require('node:stream');

const contentPolicy="default-src 'self' data: blob:; script-src 'self' 'unsafe-inline' 'wasm-unsafe-eval'; style-src 'self' 'unsafe-inline'; connect-src 'self' data: blob:; worker-src 'self' blob:; frame-src 'self'; object-src 'none'; base-uri 'none'; form-action 'none'";
const types={
  '.html':'text/html; charset=utf-8','.htm':'text/html; charset=utf-8',
  '.js':'text/javascript; charset=utf-8','.mjs':'text/javascript; charset=utf-8',
  '.css':'text/css; charset=utf-8','.json':'application/json','.txt':'text/plain; charset=utf-8',
  '.svg':'image/svg+xml','.png':'image/png','.jpg':'image/jpeg','.jpeg':'image/jpeg',
  '.webp':'image/webp','.gif':'image/gif','.avif':'image/avif','.ico':'image/x-icon',
  '.woff':'font/woff','.woff2':'font/woff2','.ttf':'font/ttf','.otf':'font/otf',
  '.wasm':'application/wasm','.mp4':'video/mp4','.webm':'video/webm',
  '.mp3':'audio/mpeg','.wav':'audio/wav','.ogg':'audio/ogg','.m4a':'audio/mp4',
};

// Stream confined local assets, including byte ranges for media seeking.
async function serveAsset(request,roots){
  let file;
  try{
    if(!['GET','HEAD'].includes(request.method))return new Response(null,{status:405});
    const url=new URL(request.url),root=Object.hasOwn(roots,url.hostname)?roots[url.hostname]:null;
    if(!root)return new Response(null,{status:404});
    const relative=decodeURIComponent(url.pathname).replace(/^\/+/,''),target=await fs.realpath(path.resolve(root,relative));
    if(target!==root&&!target.startsWith(root+path.sep))return new Response(null,{status:403});
    const info=await fs.stat(target);
    if(!info.isFile())return new Response(null,{status:404});
    file=await fs.open(target,'r');
    const {size}=await file.stat();
    const headers=new Headers({
      'Content-Type':types[path.extname(target).toLowerCase()]||'application/octet-stream',
      'Content-Security-Policy':contentPolicy,'X-Content-Type-Options':'nosniff',
      'Accept-Ranges':'bytes','Content-Length':String(size),
    });
    let start=0,end=size-1,status=200;
    const range=request.headers.get('range');
    if(range){
      const match=/^bytes=(\d*)-(\d*)$/.exec(range);
      if(match&&(match[1]||match[2])){
        start=match[1]?Number(match[1]):Math.max(0,size-Number(match[2]));
        end=match[1]&&match[2]?Math.min(Number(match[2]),size-1):size-1;
      }else start=NaN;
      if(!Number.isSafeInteger(start)||!Number.isSafeInteger(end)||start>end||start>=size){
        headers.set('Content-Range',`bytes */${size}`);headers.set('Content-Length','0');
        return new Response(null,{status:416,headers});
      }
      status=206;headers.set('Content-Range',`bytes ${start}-${end}/${size}`);
      headers.set('Content-Length',String(end-start+1));
    }
    if(request.method==='HEAD'||size===0)return new Response(null,{status,headers});
    const stream=file.createReadStream({start,end});
    file=null; // The stream owns the descriptor, including on cancellation.
    return new Response(Readable.toWeb(stream),{status,headers});
  }catch{return new Response(null,{status:404});}
  finally{await file?.close();}
}
module.exports={serveAsset};
