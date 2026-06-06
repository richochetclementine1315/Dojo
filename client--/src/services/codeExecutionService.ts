import axios from 'axios';

// Judge0 Community Edition – public instance, no API key needed for light use.
// Docs: https://ce.judge0.com/
const JUDGE0_URL = 'https://judge0-ce.p.rapidapi.com';

// RapidAPI key — leave empty to hit the open CE instance instead.
// For production, sign up at https://rapidapi.com/judge0-official/api/judge0-ce
// and put your key in the .env as VITE_JUDGE0_RAPIDAPI_KEY.
const RAPIDAPI_KEY = import.meta.env.VITE_JUDGE0_RAPIDAPI_KEY || '';

// Language IDs in Judge0 CE
// Full list: https://ce.judge0.com/languages/
const LANGUAGE_IDS: Record<string, number> = {
  javascript: 93,  // Node.js 18.15.0
  python:     71,  // Python 3.11.2
  java:       62,  // Java (OpenJDK 13.0.1)
  cpp:        54,  // C++ (GCC 9.2.0)
  c:          50,  // C (GCC 9.2.0)
  rust:       73,  // Rust (1.40.0)
  go:         60,  // Go (1.13.5)
  typescript: 74,  // TypeScript (3.7.4)
  csharp:     51,  // C# (Mono 6.6.0)
  ruby:       72,  // Ruby (2.7.0)
};

interface ExecutionResult {
  stdout: string;
  stderr: string;
  output: string;
  exitCode: number;
  time?: string;
  memory?: number;
  status?: string;
}

// Judge0 status IDs
const STATUS: Record<number, string> = {
  1:  'In Queue',
  2:  'Processing',
  3:  'Accepted',
  4:  'Wrong Answer',
  5:  'Time Limit Exceeded',
  6:  'Compilation Error',
  7:  'Runtime Error (SIGSEGV)',
  8:  'Runtime Error (SIGXFSZ)',
  9:  'Runtime Error (SIGFPE)',
  10: 'Runtime Error (SIGABRT)',
  11: 'Runtime Error (NZEC)',
  12: 'Runtime Error (Other)',
  13: 'Internal Error',
  14: 'Exec Format Error',
};

class CodeExecutionService {
  private buildHeaders() {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };
    if (RAPIDAPI_KEY) {
      headers['X-RapidAPI-Key']  = RAPIDAPI_KEY;
      headers['X-RapidAPI-Host'] = 'judge0-ce.p.rapidapi.com';
    }
    return headers;
  }

  private get baseUrl(): string {
    // If a RapidAPI key is provided use the RapidAPI endpoint,
    // otherwise fall back to the open CE instance.
    return RAPIDAPI_KEY
      ? JUDGE0_URL
      : 'https://ce.judge0.com';
  }

  async executeCode(
    language: string,
    code: string,
    stdin: string = ''
  ): Promise<ExecutionResult> {
    const languageId = LANGUAGE_IDS[language.toLowerCase()];
    if (!languageId) {
      throw new Error(`Language "${language}" is not supported`);
    }

    const encodedCode  = btoa(unescape(encodeURIComponent(code)));
    const encodedStdin = stdin ? btoa(unescape(encodeURIComponent(stdin))) : '';

    // Step 1 – submit the code
    const submitResp = await axios.post(
      `${this.baseUrl}/submissions?base64_encoded=true&wait=false`,
      {
        source_code:     encodedCode,
        language_id:     languageId,
        stdin:           encodedStdin,
        cpu_time_limit:  5,    // seconds
        memory_limit:    128000, // KB
      },
      { headers: this.buildHeaders() }
    );

    const token: string = submitResp.data.token;
    if (!token) throw new Error('No submission token received');

    // Step 2 – poll until finished (max ~15 s)
    for (let attempt = 0; attempt < 15; attempt++) {
      await new Promise(r => setTimeout(r, 1000));

      const pollResp = await axios.get(
        `${this.baseUrl}/submissions/${token}?base64_encoded=true`,
        { headers: this.buildHeaders() }
      );

      const result = pollResp.data;
      const statusId: number = result.status?.id ?? 0;

      // Still in queue / processing
      if (statusId <= 2) continue;

      const decode = (b64: string | null) =>
        b64 ? decodeURIComponent(escape(atob(b64))) : '';

      const stdout = decode(result.stdout);
      const stderr = decode(result.stderr) || decode(result.compile_output);
      const statusLabel = STATUS[statusId] ?? `Status ${statusId}`;

      let output = stdout || stderr || '';
      if (!output && statusId !== 3) output = statusLabel;

      return {
        stdout,
        stderr,
        output,
        exitCode: result.exit_code ?? (statusId === 3 ? 0 : 1),
        time:     result.time ?? undefined,
        memory:   result.memory ?? undefined,
        status:   statusLabel,
      };
    }

    throw new Error('Execution timed out waiting for Judge0 result');
  }
}

export const codeExecutionService = new CodeExecutionService();
