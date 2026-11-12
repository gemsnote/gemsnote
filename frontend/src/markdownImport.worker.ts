import {parseMarkdownImport,ImportError} from './markdownImport'
self.onmessage=({data}:{data:{name:string;bytes:Uint8Array}})=>{
 try{self.postMessage({source:parseMarkdownImport(data.name,data.bytes)})}
 catch(e){self.postMessage({error:e instanceof ImportError?e.key:e instanceof Error?e.message:'ZIP 文件损坏或格式不受支持',params:e instanceof ImportError?e.params:undefined})}
}
