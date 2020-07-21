conf = crashd_config(workdir=args.workdir)
ssh_conf = ssh_config(username="capv", private_key_path=args.private_key)
kube_conf = kube_config(path=args.mc_config)

#list out management and workload cluster nodes
wc_provider=capv_provider(
    workload_cluster=args.cluster_name,
    ssh_config=ssh_conf,
    kube_config=kube_conf
)
nodes = resources(provider=wc_provider)

capture(cmd="sudo df -i", resources=nodes)
capture(cmd="sudo crictl info", resources=nodes)
capture(cmd="df -h /var/lib/containerd", resources=nodes)
capture(cmd="sudo systemctl status kubelet", resources=nodes)
capture(cmd="sudo systemctl status containerd", resources=nodes)
capture(cmd="sudo journalctl -xeu kubelet", resources=nodes)

capture(cmd="sudo cat /var/log/cloud-init-output.log", resources=nodes)
capture(cmd="sudo cat /var/log/cloud-init.log", resources=nodes)

#add code to collect pod info from cluster
wc_kube_conf = kube_config(capi_provider = wc_provider)
set_as_default(kube_config = wc_kube_conf)

pod_ns=["default", "kube-system"]

kube_capture(what="logs", namespaces=pod_ns)
kube_capture(what="objects", kinds=["pods", "services"], namespaces=pod_ns)
kube_capture(what="objects", kinds=["deployments", "replicasets"], groups=["apps"], namespaces=pod_ns)

archive(output_file="diagnostics.tar.gz", source_paths=[conf.workdir])