const { appendFileSync, existsSync } = require('node:fs');

const requested = process.env.DISPATCH_APP;
const candidates = requested ? [requested] : JSON.parse(process.env.CHANGED_APPS || '[]');
const buildable = (app) => /^[a-z0-9][a-z0-9-]*$/.test(app)
  && existsSync(`apps/${app}/Dockerfile`)
  && existsSync(`apps/${app}/docker-bake.hcl`);

if (requested && !buildable(requested)) {
  throw new Error('The requested application has no build target');
}

const apps = [...new Set(candidates.filter(buildable))];
appendFileSync(process.env.GITHUB_OUTPUT, `changed_files=${JSON.stringify(apps)}\n`);
