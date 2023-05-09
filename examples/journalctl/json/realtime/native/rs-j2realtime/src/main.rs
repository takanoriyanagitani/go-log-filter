use std::io::{stdin, stdout, BufRead, BufReader, BufWriter, Read, Write};

use serde_json::{Map, Value};

fn read2chunks<R>(r: R) -> impl Iterator<Item = Vec<u8>>
where
    R: Read,
{
    let br = BufReader::new(r);
    let splited = br.split(b'\n');
    splited.flat_map(|r| r.ok())
}

fn chunk2realtime(c: &[u8]) -> Result<f64, String> {
    let root: Value = serde_json::from_slice(c).map_err(|e| format!("Invalid json object: {e}"))?;
    let m: &Map<String, Value> = root
        .as_object()
        .ok_or_else(|| String::from("Invalid object"))?;
    let v: &Value = m
        .get("__REALTIME_TIMESTAMP")
        .ok_or_else(|| String::from("no time"))?;
    let s: &str = v.as_str().ok_or_else(|| String::from("invalid time"))?;
    let f: f64 = str::parse(s).map_err(|e| format!("Invalid time: {e}"))?;
    Ok(f * 1e-6)
}

fn chunks2realtimes<I>(chunks: I) -> impl Iterator<Item = f64>
where
    I: Iterator<Item = Vec<u8>>,
{
    chunks.flat_map(|c: Vec<u8>| chunk2realtime(&c).ok())
}

fn compose<F, G, T, U, V>(f: F, g: G) -> impl Fn(T) -> V
where
    F: Fn(T) -> U,
    G: Fn(U) -> V,
{
    move |t: T| {
        let u: U = f(t);
        g(u)
    }
}

fn realtime2bytes(r: f64) -> [u8; 8] {
    r.to_be_bytes()
}

fn realtimes2chunks<I>(realtimes: I) -> impl Iterator<Item = [u8; 8]>
where
    I: Iterator<Item = f64>,
{
    realtimes.map(realtime2bytes)
}

fn chunk2write<W>(c: &[u8], w: &mut W) -> Result<usize, String>
where
    W: Write,
{
    w.write(c).map_err(|e| format!("Unable to write: {e}"))
}

fn chunks2write_new<W, I>(w: W) -> impl FnMut(&mut I) -> Result<usize, String>
where
    W: Write,
    I: Iterator<Item = [u8; 8]>,
{
    let mut bw = BufWriter::new(w);
    move |chunks: &mut I| {
        chunks.try_fold(0, |tot, chunk| {
            let u: &[u8] = &chunk;
            chunk2write(u, &mut bw).map(|cnt: usize| cnt + tot)
        })
    }
}

fn jsons2realtimes<R, W>(i: R, w: W) -> Result<usize, String>
where
    R: Read,
    W: Write,
{
    let realtimes = compose(read2chunks, chunks2realtimes)(i);
    let mut chunks = realtimes2chunks(realtimes);
    chunks2write_new(w)(&mut chunks)
}

fn sub() -> Result<(), String> {
    let i = stdin();
    let il = i.lock();
    let o = stdout();
    let ol = o.lock();

    jsons2realtimes(il, ol).map(|_| ())
}

fn main() {
    match sub() {
        Ok(_) => {}
        Err(e) => eprintln!("{e}"),
    }
}
