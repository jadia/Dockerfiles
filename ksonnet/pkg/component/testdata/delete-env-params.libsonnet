local params = import '../../components/params.libsonnet';

params {
  components+: {
    "nested.guestbook-ui"+: {
      name: 'guestbook-dev',
    },
  },
}