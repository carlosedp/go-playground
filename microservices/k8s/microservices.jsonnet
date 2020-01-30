local k = import 'ksonnet-libs/1.14.7/k.libsonnet';
local utils = import 'ksonnet-libs/utils.libsonnet';

local objs = {
  _config+:: {
    namespace: 'default',
    ingressSuffix: '192.168.99.100.nip.io',
    jaegerCollector: '192.168.15.141:14267',
    jaegerAgent: [
      { name: 'JAEGER_AGENT_HOST', value: 'jaeger-agent' },
      { name: 'JAEGER_AGENT_PORT', value: '6831' },
    ],
  },

  client+:: {
    deployment:
      local d = utils.newDeployment('client', $._config.namespace, 'carlosedp/microservices-demo-client', null, 8080);
      d + utils.addEnviromnentVars(d,
                                   'client',
                                   [{ name: 'JAEGER_SERVICE_NAME', value: 'client.default' }] + $._config.jaegerAgent),

    service:
      utils.newService('client', $._config.namespace, 8080),

    ingress:
      utils.newIngress('client', $._config.namespace, 'client.' + $._config.ingressSuffix, '/', 'client', 8080),
  },

  formatter+:: {
    deployment:
      local d = utils.newDeployment('formatter', $._config.namespace, 'carlosedp/microservices-demo-formatter', null, 8080);
      d + utils.addEnviromnentVars(d,
                                   'formatter',
                                   [{ name: 'JAEGER_SERVICE_NAME', value: 'formatter.default' }] + $._config.jaegerAgent)
      + utils.addContainerProperty(d, 'formatter', { imagePullPolicy: 'Always' }),

    service:
      utils.newService('formatter', $._config.namespace, 8080),
  },

  publisher+:: {
    deployment:
      local d = utils.newDeployment('publisher', $._config.namespace, 'carlosedp/microservices-demo-publisher', null, 8080);
      d + utils.addEnviromnentVars(d,
                                   'publisher',
                                   [{ name: 'JAEGER_SERVICE_NAME', value: 'publisher.default' }] + $._config.jaegerAgent),

    service:
      utils.newService('publisher', $._config.namespace, 8080),
  },

  echoapp+:: {
    deployment:
      local d = utils.newDeployment('echo-app', $._config.namespace, 'carlosedp/microservices-demo-echo-app', null, 8080);
      d + utils.addEnviromnentVars(d,
                                   'echo-app',
                                   [{ name: 'JAEGER_SERVICE_NAME', value: 'echo-app.default' }] + $._config.jaegerAgent),

    service:
      utils.newService('echo-app', $._config.namespace, 8080),
  },

  jaegeragent+:: {
    deployment:
      local d = utils.newDeployment('jaeger-agent', $._config.namespace, 'jaegertracing/jaeger-agent', null, 6831);
      d + utils.addArguments(d, 'jaeger-agent', ['--collector.host-port=jaeger-collector.istio-system.svc.cluster.local:14267']),

    service:
      local service = k.core.v1.service;
      local servicePort = k.core.v1.service.mixin.spec.portsType;
      local p = servicePort.newNamed('jaeger-agent', 6831, 6831).withProtocol('UDP');

      local s = service.new('jaeger-agent', { app: 'jaeger-agent' }, p)
                + service.mixin.metadata.withNamespace($._config.namespace)
                + service.mixin.metadata.withLabels({ app: 'jaeger-agent' });
      s,
  },
};

utils.generate(objs)
