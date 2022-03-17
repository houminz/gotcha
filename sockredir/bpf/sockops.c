// +build ignore

#include "vmlinux.h"
#include "bpf_helpers.h"
#include "bpf_endian.h"

char __license[] SEC("license") = "Dual MIT/GPL";

struct sock_key {
    __u32 sip4;
    __u32 dip4;
    __u8  family;
    __u8  pad1;   // this padding required for 64bit alignment
    __u16 pad2;   // else ebpf kernel verifier rejects loading of the program
    __u32 pad3;
    __u32 sport;
    __u32 dport;
} __attribute__((packed));


struct bpf_map_def SEC("maps") sock_ops_map = {
    .type        = BPF_MAP_TYPE_SOCKHASH,
    .key_size    = sizeof(struct sock_key),
    .value_size  = sizeof(int),
    .max_entries = 65535,
    .map_flags   = 0,
};

/*
 * extract the key identifying the socket source of the TCP event
 */
static inline
void extract_key4_from_ops(struct bpf_sock_ops *ops, struct sock_key *key)
{
    // keep ip and port in network byte order
    key->dip4 = ops->remote_ip4;
    key->sip4 = ops->local_ip4;
    key->family = 1;

    // local_port is in host byte order, and
    // remote_port is in network byte order
    key->sport = (bpf_htonl(ops->local_port) >> 16);
    key->dport = (ops->remote_port) >> 16;
    // key->dport = FORCE_READ(ops->remote_port) >> 16;
}



/*
 * Insert socket into sockmap
 */
static inline
void bpf_sock_ops_ipv4(struct bpf_sock_ops *skops)
{
    struct sock_key key = {};
    int ret;

    extract_key4_from_ops(skops, &key);

    ret = bpf_sock_hash_update(skops, &sock_ops_map, &key, BPF_NOEXIST);
    if (ret != 0) {
        bpf_trace_printk("sock_hash_update() failed, ret: %d\n", ret);
    }

    bpf_trace_printk("sockmap: op %d, port %d --> %d\n",
            skops->op, skops->local_port, bpf_ntohl(skops->remote_port));
}

/*
 * extract the key that identifies the destination socket in the sock_ops_map
 */
static inline
void extract_key4_from_msg(struct sk_msg_md *msg, struct sock_key *key)
{
    key->sip4 = msg->remote_ip4;
    key->dip4 = msg->local_ip4;
    key->family = 1;

    key->dport = (bpf_htonl(msg->local_port) >> 16);
    // key->sport = FORCE_READ(msg->remote_port) >> 16;
    key->sport = (msg->remote_port) >> 16;
}

	
SEC("sockops")
int my_sockmap(struct bpf_sock_ops *skops)
{
    switch (skops->op) {
        case BPF_SOCK_OPS_PASSIVE_ESTABLISHED_CB:
        case BPF_SOCK_OPS_ACTIVE_ESTABLISHED_CB:
            if (skops->family == 2) { //AF_INET
                bpf_sock_ops_ipv4(skops);
            }
            break;
        default:
            break;
    }
    return 0;
}


SEC("sk_msg")
int my_redir(struct sk_msg_md *msg)
{
    struct sock_key key = {};
    extract_key4_from_msg(msg, &key);
    bpf_msg_redirect_hash(msg, &sock_ops_map, &key, BPF_F_INGRESS);
    return SK_PASS;
}
