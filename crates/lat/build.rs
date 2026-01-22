fn main() {
    let target = std::env::var("TARGET").unwrap();
    if target != "wasm32-wasip2" {
        println!("cargo::warning=This crate should be built for wasm32-wasip2");
    }
}
