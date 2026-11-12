export interface NoteSelection {ids:Set<string>;anchor:string}

// Ranges follow the displayed order, including the active sort and filter.
export function selectNote(previous:NoteSelection,order:string[],id:string,modifiers:{toggle?:boolean;range?:boolean}={}):NoteSelection {
 const index=order.indexOf(id)
 if(index<0)return previous
 const visible=new Set(order),ids=new Set([...previous.ids].filter(id=>visible.has(id)))
 if(modifiers.range){
  const anchor=visible.has(previous.anchor)?previous.anchor:order[0]
  const at=order.indexOf(anchor),range=order.slice(Math.min(at,index),Math.max(at,index)+1)
  return {ids:new Set(modifiers.toggle?[...ids,...range]:range),anchor}
 }
 if(modifiers.toggle){ids.has(id)?ids.delete(id):ids.add(id);return {ids,anchor:id}}
 return {ids:new Set([id]),anchor:id}
}
