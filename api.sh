

export GENESIS_FORK_VERSION=0x10000001
export GENESIS_VALIDATORS_ROOT=0x202c996ddee3afe959f106b6759bbd3453a592f70ada65ab44b0c6dfdd0d3df5
export BELLATRIX_FORK_VERSION=0x30000001
export CAPELLA_FORK_VERSION=0x40000001
export DENEB_FORK_VERSION=0x50000001



go run . api \
  --network custom \
  --secret-key 0x607a11b45a7219cc61a3d9c5fd08c7eebd602a6a19a977f8d3771d5711a550f2 \
  --db postgres://postgres:postgres@95.217.233.186:5432/postgres?sslmode=disable \
  --beacon-uris http://116.202.172.145:5052


