# Test SCI prototype steps

```bash
# Do the normal C2P plugin setup from the QuickStart
c2pcli sci2policy -a "docs/policy.yaml"  -c docs/c2p-config.yaml --eval-dir ./testevals
c2pcli result2sci -a "./policy.yaml"  -c c2p-config.yaml --eval-dir ./testevals
c2pcli sci2posture   -a "docs/policy.yaml" --eval-dir ./testevals
```