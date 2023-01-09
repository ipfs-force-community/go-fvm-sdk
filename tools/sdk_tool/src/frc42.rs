use anyhow::Result;
use blake2b_simd::blake2b;
use clap::Parser;
use frc42_dispatch::hash::MethodResolver;
use frc42_hasher::hash::Hasher;

#[derive(Parser, Debug)]
pub struct Frc42Config {
    #[clap(last = true)]
    pub name: String,
}

pub struct Blake2bHasher {}
impl Hasher for Blake2bHasher {
    fn hash(&self, bytes: &[u8]) -> Vec<u8> {
        blake2b(bytes).as_bytes().to_vec()
    }
}

pub fn compute_frc46(cfg: &Frc42Config) -> Result<()> {
    let resolver = MethodResolver::new(Blake2bHasher {});
    let method_number = resolver.method_number(&cfg.name)?;
    println!("{}: {:#x} {}", cfg.name, method_number, method_number);
    Ok(())
}
