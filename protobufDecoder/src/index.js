import { parseInput } from "./decoder/hexUtils.js";
import { decodeProto } from "./decoder/protobufDecoder.js";
import { renderProtobufDisplay } from "./decoder/protobufDisplayServer.js";
import express from "express";
const app = express();
const port = 3000;

function iterateNestedArrays(arr) {
  var final = [];
  for (let i = 0; i < arr.length; i++) {
    if (Array.isArray(arr[i])) {
      final = [...final, ...iterateNestedArrays(arr[i])];
    } else {
      if (arr[i].type == 'Int' && i == 0 && arr[i].value.length == 7)
      {
        const value = arr[i].value;
        console.log(value);
        final.push(value);
      }
    }
  }
  return final;
}

app.use(express.json());

app.get('/', (req, res) => {
  res.send('Hello, World! I\'m Online');
});

app.post('/decode', (req, res) => {
  console.log("/decode new request")
  if (req.body.token == null || req.body.token.length == 0)
  {
    const responseData = {
      status: 400,
      message: 'Bad Request: Missing or empty request body',
    };
    res.send(responseData);
    console.error('/decode new request body invalid');
  }
  try 
  {
    const hex  = req.body.token;
    const buffer = parseInput(hex);
    const decodedResult = decodeProto(buffer);
    var decoded = renderProtobufDisplay(decodedResult);
    var result = iterateNestedArrays(decoded);
    const responseData = {
      status: 200,
      message: 'successful',
      data: result
    };
    res.send(responseData);
    console.error('/decode request ended');
  } catch (err) {
    console.error('Error:', err);
    res.send({
      status: 500,
      message: 'Internal Server Error',
    });
  }
})

app.listen(port, () => {
  console.log(`Server is listening on port ${port}`);
});
