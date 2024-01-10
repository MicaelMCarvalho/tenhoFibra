import { ProtobufPart } from "./protobufPart.js";
import { bufferToPrettyHex } from "./hexUtils.js"
var c = -1;
export function renderProtobufDisplay(value) {
  var back = value;
  const parts = value.parts.map((part, i) => {
    return ProtobufPart(part); 
  });

  const leftOver = value.leftOver.length ? {
    leftOver: bufferToPrettyHex(value.leftOver),
  } : null;

  return parts;
}
