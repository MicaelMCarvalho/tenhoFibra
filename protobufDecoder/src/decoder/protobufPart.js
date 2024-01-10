import { decodeProto, TYPES, typeToString } from "./protobufDecoder.js";
import {
  decodeFixed32,
  decodeFixed64,
  decodeStringOrBytes,
  decodeVarintParts
} from "./protobufPartDecoder.js";
import { renderProtobufDisplay } from "./protobufDisplayServer.js";

function ProtobufVarintPart(value) {
  const decoded = decodeVarintParts(value);

  return decoded.map((d, i) => ({
    type: d.type,
    value: d.value
  }));
}

function ProtobufStringOrBytesPart(value) {
  return { value: value.value };
}

function ProtobufFixed64Part(value) {
  const decoded = decodeFixed64(value);
  return decoded.map((d, i) => ({
    type: d.type,
    value: d.value
  }));
}

function ProtobufFixed32Part(value) {
  const decoded = decodeFixed32(value);
  return decoded.map((d, i) => ({
    type: d.type,
    value: d.value
  }));
}

function getProtobufPart(part) {
  switch (part.type) {
    case TYPES.VARINT:
      return ProtobufVarintPart(part.value);
    case TYPES.LENDELIM:
      let decoded = decodeProto(part.value);
      if (part.value.length > 0 && decoded.leftOver.length === 0) {
        return [renderProtobufDisplay(decoded), "protobuf"];
      } else {
        decoded = decodeStringOrBytes(part.value);
        return [ProtobufStringOrBytesPart(decoded), decoded.type] ;
      }
    case TYPES.FIXED64:
      var a = ProtobufFixed64Part(part.value);
      return a;
    case TYPES.FIXED32:
      var b = ProtobufFixed32Part(part.value);
      return b;
    default:
      return { type: "Unknown type" };
  }
}

export function ProtobufPart(part) {
  const { byteRange, index, type } = part;
  var contents = getProtobufPart(part);
  return contents;
}
