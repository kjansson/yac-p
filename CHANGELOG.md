# Changelog

## [3.1.1](https://github.com/kjansson/yac-p/compare/v3.1.0...v3.1.1) (2025-06-18)


### Bug Fixes

* update release-please config to trigger patch release on deps update ([#56](https://github.com/kjansson/yac-p/issues/56)) ([391c5e5](https://github.com/kjansson/yac-p/commit/391c5e5e0bdec68e07c45f4a8eb927c2c39c056a))

## [3.1.0](https://github.com/kjansson/yac-p/compare/v3.0.1...v3.1.0) (2025-05-13)


### Features

* error checks ([#53](https://github.com/kjansson/yac-p/issues/53)) ([c8e9095](https://github.com/kjansson/yac-p/commit/c8e9095d324abfb38bba400f0bd136df9a0896d9))

## [3.0.1](https://github.com/kjansson/yac-p/compare/v3.0.0...v3.0.1) (2025-04-30)


### Bug Fixes

* comments for correct package docs ([#50](https://github.com/kjansson/yac-p/issues/50)) ([1e9ad7b](https://github.com/kjansson/yac-p/commit/1e9ad7bd0596d389bb3fd0c4f0c4ac32c8c7a05d))

## [3.0.0](https://github.com/kjansson/yac-p/compare/v2.0.0...v3.0.0) (2025-04-30)


### ⚠ BREAKING CHANGES

* new interface to interact with controller and all components

### Features

* remove inter-dependencies ([#44](https://github.com/kjansson/yac-p/issues/44)) ([02ef8af](https://github.com/kjansson/yac-p/commit/02ef8afaacaef69233f0da52dc957eae883c0eb6))
* remove yace hard dependencies ([#48](https://github.com/kjansson/yac-p/issues/48)) ([c0c9bf5](https://github.com/kjansson/yac-p/commit/c0c9bf55a33683791ca029352e94ed0d7b4bca72))

## [2.0.0](https://github.com/kjansson/yac-p/compare/v1.0.0...v2.0.0) (2025-04-23)


### ⚠ BREAKING CHANGES

* Moves all env var configuration out of packages, config now part of types

### Features

* reorganize packages ([#37](https://github.com/kjansson/yac-p/issues/37)) ([46878b8](https://github.com/kjansson/yac-p/commit/46878b8a5825642a750e15d41a6a161c64e89aa4))

## [1.0.0](https://github.com/kjansson/yac-p/compare/v0.9.0...v1.0.0) (2025-04-22)


### ⚠ BREAKING CHANGES

* changes field names in the controller struct

### Features

* code makeover ([#34](https://github.com/kjansson/yac-p/issues/34)) ([5653a84](https://github.com/kjansson/yac-p/commit/5653a847051a16b4bdb3a7e74470f743fe25471f))

## [0.9.0](https://github.com/kjansson/yac-p/compare/v0.8.1...v0.9.0) (2025-04-22)


### Features

* adding debug logging ([#32](https://github.com/kjansson/yac-p/issues/32)) ([7c12f00](https://github.com/kjansson/yac-p/commit/7c12f00682f2bd318934704c1824ade13e8d553f))

## [0.8.1](https://github.com/kjansson/yac-p/compare/v0.8.0...v0.8.1) (2025-04-22)


### Bug Fixes

* cleanup ([#28](https://github.com/kjansson/yac-p/issues/28)) ([16670c9](https://github.com/kjansson/yac-p/commit/16670c936d24d1908d1dba0da6b4c1fddc0fc7b7))

## [0.8.0](https://github.com/kjansson/yac-p/compare/v0.7.1...v0.8.0) (2025-04-21)


### Features

* remove hard dependencies and add tests ([#26](https://github.com/kjansson/yac-p/issues/26)) ([4c4e41d](https://github.com/kjansson/yac-p/commit/4c4e41d12742a0b75b2f9c35a57612d08aaa0614))


### Bug Fixes

* remove the need for remote write role in assumable roles ([#19](https://github.com/kjansson/yac-p/issues/19)) ([17b9c19](https://github.com/kjansson/yac-p/commit/17b9c197b27edf3fd60a479d5c520b2df0da1f01))
* rename types file for tests ([#27](https://github.com/kjansson/yac-p/issues/27)) ([7e97fc9](https://github.com/kjansson/yac-p/commit/7e97fc9192159c82f343258ef719c0eba16be622))

## [0.7.1](https://github.com/kjansson/yac-p/compare/v0.7.0...v0.7.1) (2025-03-31)


### Bug Fixes

* debug logging ([#17](https://github.com/kjansson/yac-p/issues/17)) ([16cb9fc](https://github.com/kjansson/yac-p/commit/16cb9fc3bc625dd73b265810d0921118a31d11a3))

## [0.7.0](https://github.com/kjansson/yac-p/compare/v0.6.0...v0.7.0) (2025-03-31)


### Features

* update to 0.62.1 libraries ([#15](https://github.com/kjansson/yac-p/issues/15)) ([df66158](https://github.com/kjansson/yac-p/commit/df66158a1887266063904372c4c23de431b27d48))

## [0.6.0](https://github.com/kjansson/yac-p/compare/v0.5.1...v0.6.0) (2025-03-25)


### Features

* move from using ecr image to al2 runtime ([#12](https://github.com/kjansson/yac-p/issues/12)) ([29e287a](https://github.com/kjansson/yac-p/commit/29e287aa9d357eb7156338b7b7bd32d2153c36f5))

## [0.5.1](https://github.com/kjansson/yac-p/compare/v0.5.0...v0.5.1) (2025-03-25)


### Bug Fixes

* remove need for amp region in config ([59ef67f](https://github.com/kjansson/yac-p/commit/59ef67fe1a349ccfe684b97c35a26fa97c6ab9b0))
* tighten down s3 policy ([d830740](https://github.com/kjansson/yac-p/commit/d8307400a206f8806a91a516ddcba77fba712292))

## [0.5.0](https://github.com/kjansson/yac-p/compare/v0.4.0...v0.5.0) (2025-03-25)


### Features

* configurable remote write role ([10b26e7](https://github.com/kjansson/yac-p/commit/10b26e7315c570ffa08cf2f03668e64b5be7e760))


### Bug Fixes

* add correct ecr permissions ([f85060e](https://github.com/kjansson/yac-p/commit/f85060e05d799bff64fbf1b7da89e238be31efca))
* condition on amp endpoint output ([672e73c](https://github.com/kjansson/yac-p/commit/672e73c02353de9a09af61dd908ed14cdde30e65))

## [0.4.0](https://github.com/kjansson/yac-p/compare/v0.3.1...v0.4.0) (2025-03-25)


### Features

* add variable controlling assumable roles ([#8](https://github.com/kjansson/yac-p/issues/8)) ([b27cdeb](https://github.com/kjansson/yac-p/commit/b27cdeb3f5f7998d79d64a6d38b0abc23886f738))

## [0.3.1](https://github.com/kjansson/yac-p/compare/v0.3.0...v0.3.1) (2025-03-24)


### Bug Fixes

* make schedule expression configurable through variable ([d21cbcb](https://github.com/kjansson/yac-p/commit/d21cbcbbd0c65feda1a5845d45a1c462678a6a8d))

## [0.3.0](https://github.com/kjansson/yac-p/compare/v0.2.0...v0.3.0) (2025-03-23)


### Features

* enable internal metrics ([#5](https://github.com/kjansson/yac-p/issues/5)) ([6a54f8e](https://github.com/kjansson/yac-p/commit/6a54f8efb666c20a0a58653f8700ea23bc84962f))

## [0.2.0](https://github.com/kjansson/yac-p/compare/v0.1.0...v0.2.0) (2025-03-23)


### Features

* extend concurrency configuration ([c4c8ae2](https://github.com/kjansson/yac-p/commit/c4c8ae20ea79fed851b7558330867f02bc9b12d8))

## [0.1.0](https://github.com/kjansson/yac-p/compare/v0.0.2...v0.1.0) (2025-03-22)


### Features

* added release please ([3748389](https://github.com/kjansson/yac-p/commit/374838910b3f32422eb5ea902709cc510249e601))
