local k = import '1.14.7/k.libsonnet';
local utils = import 'utils.libsonnet';

local objs = {
  _config+:: {
    appname: 'nginx',
    namespace: 'default',
    port: 80,
    image: 'nginx:1.17',
    ingressHost: '192.168.99.102.nip.io'
  },

  demoApp+:: {
    deployment:
      local d = utils.newDeployment($._config.appname, $._config.namespace, $._config.image, null, $._config.port);
      d + { spec+: { replicas: 3 }}, // Change default number of replicas

    service:
      utils.newService($._config.appname, $._config.namespace, $._config.port),

    ingress:
      utils.newIngress($._config.appname, $._config.namespace, $._config.appname+'.'+$._config.ingressHost, '/', $._config.appname, $._config.port)
  },
};

utils.generate(objs)