package pixie

var podStatsQuery = `import px

df = px.DataFrame(table='process_stats', start_time='-60s')
df.container = df.ctx['container']
df.pod = df.ctx['pod']
df.pid = df.ctx['pid']
df.node = df.ctx['node']
df.cmd = df.ctx['cmd']
df.namespace = df.ctx['namespace']
df.rss = df.rss_bytes
df.vsz = df.vsize_bytes
df.time = df.cpu_utime_ns + df.cpu_ktime_ns

df = df.groupby(['pid', 'container', 'pod', 'rss', 'vsz', 'node', "namespace", "cmd"]).agg(
    time=('time', px.max),
)

px.display(df)`
