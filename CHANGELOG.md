# Changelog

## [1.3.1](https://github.com/albugowy15/api-double-track/compare/v1.3.0...v1.3.1) (2024-07-12)


### Bug Fixes

* typo in alternatives description ([ffceeee](https://github.com/albugowy15/api-double-track/commit/ffceeee092a5cdaf49a02e7de2fa180f8bdf1608))

## [1.3.0](https://github.com/albugowy15/api-double-track/compare/v1.2.0...v1.3.0) (2024-06-27)


### Features

* add CalculateAlternativeHpt in ahp package ([2270bb6](https://github.com/albugowy15/api-double-track/commit/2270bb627f41d4f249b796352a69c1caf7728e61))


### Bug Fixes

* typo alternative "Teknik Kendaraan Ringan/Motor" ([00c75a8](https://github.com/albugowy15/api-double-track/commit/00c75a886a251ab08883687f5f06099796f1acb2))

## [1.2.0](https://github.com/albugowy15/api-double-track/compare/v1.1.0...v1.2.0) (2024-06-19)


### Features

* use jwtauth v5 and jwx v2 ([7a929d2](https://github.com/albugowy15/api-double-track/commit/7a929d2b5cf1e7d9fd1d3c51c910ea748b6428b2))

## [1.1.0](https://github.com/albugowy15/api-double-track/compare/v1.0.0...v1.1.0) (2024-06-19)


### Features

* add ahp weight in topsis method ([da938ab](https://github.com/albugowy15/api-double-track/commit/da938abd117fa4dec737a3130655987230b6ac5a))
* add ahp weight in topsis method ([07912de](https://github.com/albugowy15/api-double-track/commit/07912dec465c77227ff24251df3699b316756557))
* add combinative weight (AHP & ENTROPY) ([551f167](https://github.com/albugowy15/api-double-track/commit/551f167144f5c67d73a461b44c7ac3e31b3c0b18))
* add entropy and topsis ([9ac9dbf](https://github.com/albugowy15/api-double-track/commit/9ac9dbfdbf6835c436e14b24284455162d0df18f))
* add expectations ([caf64d5](https://github.com/albugowy15/api-double-track/commit/caf64d54cdc5307bd1963e277f191d10fcf57878))
* add expectations and expectations_to_alternative migrations ([ce7a476](https://github.com/albugowy15/api-double-track/commit/ce7a476a7eb930b9a4ba2d3fde44f2fea73e2ee4))
* add ON DELETE CASCADE ([cbe338b](https://github.com/albugowy15/api-double-track/commit/cbe338b9d54c1c9eaf65e0e4b7607d1cbffc8b34))
* add student register endpoint ([da04230](https://github.com/albugowy15/api-double-track/commit/da04230731daa482c50cacd719df1b0bee88d7d7))
* admin and student change password route ([ba211e0](https://github.com/albugowy15/api-double-track/commit/ba211e003eea762f5d75d57471f135a740863e7c))
* make multiple recommendation items can have same rank ([98d9152](https://github.com/albugowy15/api-double-track/commit/98d9152299cf4f5850e91998fbca6631d87b03a5))
* move calculateCriteriaWeight to controllers ([bd72267](https://github.com/albugowy15/api-double-track/commit/bd722672c0331c9a5a40669b5d8d42077915fc46))
* order schools by name asc ([7895d2a](https://github.com/albugowy15/api-double-track/commit/7895d2a3bf64a746106668647ada861bc208fbdd))
* topsis recommendations by school query ([51fa5f7](https://github.com/albugowy15/api-double-track/commit/51fa5f72e0e3868d29c84876a7ffb978c9822b4c))


### Bug Fixes

* add student id foreign key for ahp and topsis ([f0fb31b](https://github.com/albugowy15/api-double-track/commit/f0fb31b199fea372a7c71bc8a896c99b886a9ed2))
* admin reccomendation ([25c6dd1](https://github.com/albugowy15/api-double-track/commit/25c6dd1282fa6c55b5828d7aeb850cdee73c6341))
* check unique field when updating student ([a32a81d](https://github.com/albugowy15/api-double-track/commit/a32a81dcdfb252e756219e4cb4f224f6650ffc58))
* database topsis_combinatives ([bc5a807](https://github.com/albugowy15/api-double-track/commit/bc5a807ea0549e0bc3a7209e4180b01c1ef94ce4))
* database topsis_combinatives ([3d1e02e](https://github.com/albugowy15/api-double-track/commit/3d1e02e6a75d9d4442cc2184bc71c681b26b6d4f))
* entropy topsis ([b50d276](https://github.com/albugowy15/api-double-track/commit/b50d2765be0b79175dcb8d0dbd007fbc6b625725))
* get recommendations by school query ([f6fa34f](https://github.com/albugowy15/api-double-track/commit/f6fa34f3603ed3aecda4d08326df7a24b20695b2))
* ignore password length when login and change password (oldPassword) ([f89ba59](https://github.com/albugowy15/api-double-track/commit/f89ba5983742cf6980543ec643d440501c24aa37))
* phone number prefix validation ([1f52fb4](https://github.com/albugowy15/api-double-track/commit/1f52fb4806c97bce5bf2bfa5053cfe128c681df1))
* typo migrate down expectations_to_alternatives ([2753480](https://github.com/albugowy15/api-double-track/commit/27534807fd7142fde6e103f03b7dc5582815ee59))
* username and password validation between login and register ([88b881a](https://github.com/albugowy15/api-double-track/commit/88b881af8b0d185f1a33d2f9519c5ddf8cd90442))
