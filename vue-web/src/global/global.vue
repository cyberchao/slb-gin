<script>
const nginx_clusters = [
  {
    value: "nginx-ssl",
    type: "success",
  },
  {
    value: "nginx-bc",
    type: "",
  },
  {
    value: "nginx-cproxy",
    type: "warning",
  },
  {
    value: "pafm-fs",
    type: "danger",
  },
];
const envs = [
  {
    value: "uat",
  },
  {
    value: "hd",
  },
  {
    value: "prd_bx",
  },
  {
    value: "prd_wgq",
  },
];
export default {
  nginx_clusters: nginx_clusters,
  envs: envs,
};
</script>
