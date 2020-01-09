'use strict';

import { ConsoleSpanExporter, SimpleSpanProcessor } from '@opentelemetry/tracing';
import { WebTracer } from '@opentelemetry/web';
import { XMLHttpRequestPlugin } from '@opentelemetry/plugin-xml-http-request';
import { ZoneScopeManager } from '@opentelemetry/scope-zone';

const webTracerWithZone = new WebTracer({
  scopeManager: new ZoneScopeManager(),
  plugins: [
    new XMLHttpRequestPlugin({
      ignoreUrls: [/localhost:8090\/sockjs-node/],
      propagateTraceHeaderCorsUrls: [
        'http://localhost:3000/calculate'
      ]
    })
  ]
});

webTracerWithZone.addSpanProcessor(new SimpleSpanProcessor(new ConsoleSpanExporter()));


const handleForm = () => {
    const endpoint = 'http://localhost:3000/calculate'
    let form = document.getElementById('calc')

    const onClick = (event) => {
        event.preventDefault();
        const span = webTracerWithZone.startSpan('calc-request', { parent: webTracerWithZone.getCurrentSpan() });
        let fd = new FormData(form);
        let requestPayload = {
            method: fd.get('calcMethod'),
            operands: tokenizeOperands(fd.get('values'))
        };
        webTracerWithZone.withSpan(span, () => {
          calculate(endpoint, requestPayload).then((res) => {
            webTracerWithZone.getCurrentSpan().addEvent('request-complete');
            span.end();
            updateResult(res);
          });
        });
    }
    form.addEventListener('submit', onClick)
}

const calculate = (endpoint, payload) => {
  return new Promise(async (resolve, reject) => {
    const req = new XMLHttpRequest();
    req.open('POST', endpoint, true);
    req.setRequestHeader('Content-Type', 'application/json');
    req.setRequestHeader('Accept', 'application/json');
    req.send(JSON.stringify(payload))
    req.onload = function () {
      resolve(req.response);
    };
  });
};

const updateResult = (res) => {
  document.getElementById('result').innerHTML = res
}

const tokenizeOperands = (values) => {
  return values.split(',').map(Number)
}

window.addEventListener('load', handleForm);
