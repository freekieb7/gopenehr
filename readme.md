# Flags

| Name | Type | Default |
| ----- | ----- | ------ |
| `prod` | Boolean | `true` |
| `port` | String | `:3000` |

# ArgoCD central repo

```
source:
  repoURL: https://github.com/freekieb7/gopenehr
  path: charts/gopenehr
  targetRevision: master
```