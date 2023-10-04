# Changelog

## 0.3.0

- Renamed "marco" to "service call" for simplification (breaking change!)
- Added multiarch support: armv7 and aarch64 (arm64) are now supported!
- Added some smaller internal frontend improvements

## 0.2.1

- Feature: Comment is now shown on the submission notification if available

## 0.2.0

- Feature: Tokens can now be used multiple times
- This release is breaking as the database format changed!

## 0.1.3

- UI: Removed Token Modal - displaying all information and delete action in main view
- Easier Token copying, as clipboard API cannot be used in Webview (Home Assistant Companion) 

## 0.1.2

- Bugfix: Macro state when creating a token was sometimes not in sync with actual selection
- Bumped Go from 1.18 to 1.20
- New Feature: Implemented notifications on when a token is submitted 
- Added tzdata to image for correct timezone info

## 0.1.1

- Fixed wrong check of development flags, preventing the API server to enter production mode
- Added absolute paths to README.md

## 0.1.0

- Initial version
