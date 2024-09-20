#![cfg_attr(target_arch = "wasm32", no_std)]
extern crate alloc;
extern crate core;

// #[cfg(feature = "blended")]
// mod blended;
#[cfg(feature = "evm")]
mod evm;
#[cfg(feature = "fvm")]
mod fvm;
#[cfg(any(
    feature = "blake2",
    feature = "sha256",
    feature = "ripemd160",
    feature = "identity",
    feature = "modexp",
    feature = "ecrecover",
))]
mod precompile;
#[cfg(feature = "svm")]
mod svm;
#[cfg(feature = "wasm")]
mod wasm;
