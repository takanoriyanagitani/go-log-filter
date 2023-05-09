static _INPUT_JSON: [u8; 65536] = [0; 65536];

#[derive(serde::Deserialize)]
struct RawLog {
    #[serde(rename(deserialize = "__REALTIME"))]
    __realtime: String,
}

struct Log {
    __realtime: f64,
}

impl TryFrom<&[u8]> for RawLog {
    type Error = String;
    fn try_from(s: &[u8]) -> Result<Self, Self::Error> {
        serde_json::from_slice(s).map_err(|e| format!("Unable to convert: {e}"))
    }
}

impl TryFrom<&RawLog> for Log {
    type Error = String;
    fn try_from(r: &RawLog) -> Result<Self, Self::Error> {
        let s: &str = r.__realtime.as_str();
        let i: i64 = str::parse(s).map_err(|e| format!("Invalid realtime: {e}"))?;
        let f: f64 = (i as f64) * 1e-6;
        Ok(Log { __realtime: f })
    }
}

fn _input_json() -> &'static [u8] {
    _INPUT_JSON.as_ref()
}

fn _json2log(s: &[u8]) -> Result<Log, String> {
    let r: RawLog = RawLog::try_from(s)?;
    Log::try_from(&r)
}

fn _json2realtime(s: &[u8]) -> Result<f64, String> {
    let l: Log = _json2log(s)?;
    Ok(l.__realtime)
}

#[no_mangle]
pub extern "C" fn to_real() -> f64 {
    _json2realtime(_input_json()).ok().unwrap_or(f64::NAN)
}
