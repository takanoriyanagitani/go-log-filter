#![deny(unsafe_code)]

use serde_json::{Map, Value};

static _INPUT_JSON: [u8; 65536] = [0; 65536];

static mut _EMSG: [u8; 65536] = [0; 65536];
static mut _ELEN: u16 = 0;

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn eaddr() -> *const u8 {
    unsafe { _EMSG.as_ptr() }
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn esize() -> u16 {
    unsafe { _ELEN }
}

fn _eaddr() -> &'static mut [u8] {
    #[allow(unsafe_code)]
    unsafe {
        _EMSG.as_mut()
    }
}

fn _ecopy(dst: &mut [u8], src: &[u8]) {
    let max: usize = src.len() & 0xffff;
    let limited: &[u8] = &src[..max];
    let l: usize = limited.len();
    let target: &mut [u8] = &mut dst[..l];
    target.copy_from_slice(limited);
    #[allow(unsafe_code)]
    unsafe {
        _ELEN = target.len() as u16
    };
}

fn _error(emsg: &str) {
    let eb: &[u8] = emsg.as_bytes();
    let dst: &mut [u8] = _eaddr();
    _ecopy(dst, eb)
}

pub struct RawLog {
    realtime: String,
}

struct Log {
    realtime: f64,
}

impl TryFrom<&[u8]> for RawLog {
    type Error = String;
    fn try_from(s: &[u8]) -> Result<Self, Self::Error> {
        let root: Value =
            serde_json::from_slice(s).map_err(|e| format!("Unable to convert to a value: {e}"))?;
        let m: &Map<String, Value> = root
            .as_object()
            .ok_or_else(|| String::from("Invalid input"))?;
        let v: &Value = m
            .get("__REALTIME_TIMESTAMP")
            .ok_or_else(|| String::from("no time"))?;
        let s: &str = v.as_str().ok_or_else(|| String::from("Invalid time"))?;
        Ok(Self { realtime: s.into() })
    }
}

impl TryFrom<&RawLog> for Log {
    type Error = String;
    fn try_from(r: &RawLog) -> Result<Self, Self::Error> {
        let s: &str = r.realtime.as_str();
        let i: i64 = str::parse(s).map_err(|e| format!("Invalid realtime: {e}"))?;
        let f: f64 = (i as f64) * 1e-6;
        Ok(Log { realtime: f })
    }
}

fn _input_json(len: u16) -> &'static [u8] {
    let all: &[u8] = _INPUT_JSON.as_ref();
    let l: usize = len as usize;
    &all[..l]
}

fn _json2log(s: &[u8]) -> Result<Log, String> {
    let r: RawLog = RawLog::try_from(s)?;
    Log::try_from(&r)
}

fn _json2realtime(s: &[u8]) -> Result<f64, String> {
    let l: Log = _json2log(s)?;
    Ok(l.realtime)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn addr() -> *const u8 {
    _INPUT_JSON.as_ptr()
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn to_real(len: u16) -> f64 {
    _json2realtime(_input_json(len)).ok().unwrap_or(f64::NAN)
}
