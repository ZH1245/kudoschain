# kudoschain

A Cosmos SDK blockchain scaffolded with [Ignite CLI](https://ignite.com/cli), extended with **two custom modules** — `kudos` and `bounty`.

Built while learning Cosmos: the point was to implement the full module stack end to end — protos → store keys → keeper → MsgServer → queries → CLI/REST — and to understand the transaction lifecycle, module accounts, escrow and per-block logic.

- **Cosmos SDK** v0.53.6 · **CometBFT** v0.38.21 · **Go** 1.26+
- Binary: `kudoschaind` · chain id: `kudoschain`
- Scaffolded with Ignite CLI; no dependency on any other chain.

---

## Modules

### `kudos` — give kudos
A small social module: record a kudos (giver → receiver + message) on-chain. No tokens involved.

| | |
|---|---|
| **Tx** `MsgGiveKudos` → `give-kudos` | `creator`, `receiver`, `message` → returns the new `id`. Validates both addresses; ID comes from a `collections.Sequence`. |
| **Query** `ShowKudo` → `show-kudo <id>` | returns `Kudo { giver, receiver, message }`. |
| **Query** `Params` | module parameters. |
| **Storage** | `Map[uint64, Kudo]` (`kudos_list`), sequence (`kudos_sq`), params (`p_kudos`) |

### `bounty` — escrowed task with expiry refund
A bounty **escrows** a reward from the creator into the module account, pays a claimer within its lifetime, and **refunds the creator automatically** if it is never claimed. This mirrors Lumera's "action fee escrow + expiry refund".

| | |
|---|---|
| **Tx** `MsgCreateBounty` → `create-bounty` | `title`, `reward` (coin) → escrows `reward` from the creator to the `bounty` **module account** (`SendCoinsFromAccountToModule`); stores the bounty with `expires_at = blockHeight + bounty_lifetime_blocks`. Emits `bounty_created`. |
| **Tx** `MsgClaimBounty` → `claim-bounty <id>` | pays the reward from the module account to the claimer (`SendCoinsFromModuleToAccount`), marks it claimed. Emits `bounty_claimed`. Errors: **`1101` already claimed**, **`1102` expired**, **`5` not found**. |
| **Query** `ShowBounty` → `show-bounty <id>` | returns `Bounty { title, reward, creator, claimed, expires_at }`. |
| **EndBlock** | `RefundExpiredBounties` walks unclaimed bounties past `expires_at`, refunds the creator and deletes the record. Emits `bounty_refunded`. |
| **Params** | `bounty_lifetime_blocks` (default **20**). |
| **Storage** | `Map[uint64, Bounty]` (`bounty_list`), sequence (`bounty_seq`), params (`p_bounty`) |

---

## Run it

```bash
ignite chain serve            # builds, initializes and starts the chain
                              # after changing code, restart with: ignite chain serve --reset-once
```
Ports: RPC `26657`, REST `1317`, gRPC `9090`. Chain id `kudoschain`.

### Demo

```bash
# kudos (no tokens)
kudoschaind tx kudos give-kudos $(kudoschaind keys show bob -a --keyring-backend test) "nice work" \
  --from alice --chain-id kudoschain --keyring-backend test --yes
kudoschaind q kudos show-kudo 0

# bounty: escrow -> claim
kudoschaind tx bounty create-bounty "Fix the README" 500token \
  --from alice --chain-id kudoschain --keyring-backend test --yes
kudoschaind q bounty show-bounty 0
kudoschaind tx bounty claim-bounty 0 --from bob --chain-id kudoschain --keyring-backend test --yes
kudoschaind q bank balances $(kudoschaind keys show bob -a --keyring-backend test)   # bob +500token

# a second claim fails with code 1101;
# a bounty left unclaimed past expires_at is refunded automatically in EndBlock (event bounty_refunded)
```

REST equivalents:
```
GET http://localhost:1317/kudoschain/kudos/v1/show_kudo/0
GET http://localhost:1317/kudoschain/bounty/v1/show_bounty/0
```

---

## Where things live

| Piece | Path |
|---|---|
| Protos (message / query / state shapes) | `proto/kudoschain/{kudos,bounty}/v1/*.proto` |
| Keepers (storage + logic) | `x/kudos/keeper/`, `x/bounty/keeper/` |
| MsgServers | `x/*/keeper/msg_server_*.go` |
| Queries | `x/*/keeper/query_*.go` |
| Errors (codes 1100–1102) | `x/bounty/types/errors.go` |
| EndBlock refund | `x/bounty/keeper/expire.go`, wired in `x/bounty/module/module.go` |
| Store keys | `x/*/types/keys.go` |
| Expected keepers (bank/auth) | `x/bounty/types/expected_keepers.go` |

---

## What this demonstrates

- **MsgServer = controller, Keeper = repository.** Only the keeper touches the module store; `collections` prefixes act as table names.
- **Escrow pattern.** `SendCoinsFromAccountToModule` on create, `SendCoinsFromModuleToAccount` on claim/refund — the same shape as Lumera's action-fee escrow and payout.
- **Per-block logic.** `EndBlock` runs every block with no transaction (a "cron job"): the automatic refund on expiry.
- **Determinism.** Handlers use only block height/time — no wall-clock, no randomness; every validator must compute the same result.
- **Two different counters.** A per-module ID `Sequence` vs the built-in per-account sequence that prevents replay attacks.
