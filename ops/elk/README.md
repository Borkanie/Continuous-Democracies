ELK (Elasticsearch + Kibana) + Filebeat (self-hosted)
===============================================

This is a lightweight local/dev ELK setup to try out logs for this repo.

What it contains
- Elasticsearch (single-node)
- Kibana
- Filebeat (reads Docker container logs and forwards to Elasticsearch)

Quick start
1. From the repo root run:

```bash
cd ops/elk
docker-compose up --build
```

2. Wait for Elasticsearch and Kibana to be ready.
   - Elasticsearch: http://localhost:9200
   - Kibana: http://localhost:5601

3. To load index templates / ILM in newer Filebeat, run filebeat setup (optional):

```bash
docker exec -it filebeat filebeat setup --index-management --pipelines --dashboards
```

Notes and limitations
- This is intended for local testing only. It disables security on Elasticsearch/Kibana.
- Filebeat is configured to read Docker container logs from `/var/lib/docker/containers`. Ensure your runner host uses this layout.
- ILM is enabled with a rollover alias; the default retention is 10 days by ILM policy (you must apply or customize the policy in Elasticsearch).

Next steps
- If you want, I can add a small ILM policy file and an example `filebeat setup` job to bootstrap the index lifecycle and templates automatically.
