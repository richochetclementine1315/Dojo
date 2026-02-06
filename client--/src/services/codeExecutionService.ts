import axios from 'axios';

const PISTON_API_URL = 'https://emkc.org/api/v2/piston';

interface ExecutionResult {
  stdout: string;
  stderr: string;
  output: string;
  exitCode: number;
}

interface PistonRuntime {
  language: string;
  version: string;
}

// Language mapping from our UI to Piston API
const languageMap: Record<string, string> = {
  javascript: 'javascript',
  python: 'python',
  java: 'java',
  cpp: 'c++',
};

class CodeExecutionService {
  async getRuntimes(): Promise<PistonRuntime[]> {
    try {
      const response = await axios.get(`${PISTON_API_URL}/runtimes`);
      return response.data;
    } catch (error) {
      throw new Error('Failed to fetch available runtimes');
    }
  }

  async executeCode(
    language: string,
    code: string,
    stdin: string = ''
  ): Promise<ExecutionResult> {
    try {
      const pistonLanguage = languageMap[language] || language;
      
      const response = await axios.post(`${PISTON_API_URL}/execute`, {
        language: pistonLanguage,
        version: '*', // Use latest version
        files: [
          {
            name: this.getFileName(language),
            content: code,
          },
        ],
        stdin: stdin,
        args: [],
        compile_timeout: 10000,
        run_timeout: 3000,
        compile_memory_limit: -1,
        run_memory_limit: -1,
      });

      const { run, compile } = response.data;
      
      // Combine compilation and runtime output
      let output = '';
      let stderr = '';
      let stdout = '';

      if (compile) {
        if (compile.stdout) stdout += compile.stdout;
        if (compile.stderr) stderr += compile.stderr;
        output += compile.output || '';
      }

      if (run) {
        if (run.stdout) stdout += run.stdout;
        if (run.stderr) stderr += run.stderr;
        output += run.output || '';
      }

      return {
        stdout,
        stderr,
        output: output || stdout || stderr,
        exitCode: run?.code || 0,
      };
    } catch (error: any) {
      throw new Error(
        error.response?.data?.message || 'Failed to execute code'
      );
    }
  }

  private getFileName(language: string): string {
    const extensions: Record<string, string> = {
      javascript: 'script.js',
      python: 'script.py',
      java: 'Main.java',
      cpp: 'main.cpp',
    };
    return extensions[language] || 'script.txt';
  }
}

export const codeExecutionService = new CodeExecutionService();
