const ModuleFederationPlugin = require("webpack/lib/container/ModuleFederationPlugin");
const mf = require("@angular-architects/module-federation/webpack");
const path = require("path");
const share = mf.share;

const sharedMappings = new mf.SharedMappings();
sharedMappings.register(
  path.join(__dirname, '../../tsconfig.json'),
  [
    'auth-lib',
    'data-lib'
  ]);

module.exports = {
  output: {
    uniqueName: "shell",
    publicPath: "auto"
  },
  optimization: {
    runtimeChunk: false
  },
  resolve: {
    alias: {
      ...sharedMappings.getAliases(),
    }
  },
  plugins: [
    new ModuleFederationPlugin({

        // For remotes (please adjust)
        // name: "shell",
        // filename: "remoteEntry.js",
        // exposes: {
        //     './Component': './projects/shell/src/app/app.component.ts',
        // },

        // For hosts (please adjust)
        remotes: { },

        shared: share({
          "@angular/core": { singleton: true, strictVersion: true, requiredVersion: 'auto' },
          "@angular/common": { singleton: true, strictVersion: true, requiredVersion: 'auto' },
          "@angular/common/http": { singleton: true, strictVersion: true, requiredVersion: 'auto' },
          "@angular/router": { singleton: true, strictVersion: true, requiredVersion: 'auto' },
          "keycloak-angular": { singleton: true, strictVersion: true, requiredVersion: 'auto' },
          "@apollo/client": { singleton: true, strictVersion: true, requiredVersion: 'auto' },
          "apollo-angular": { singleton: true, strictVersion: true, requiredVersion: 'auto' },
          "@angular/material/input": { singleton: true, strictVersion: true, requiredVersion: 'auto'},
          "@angular/material/select": { singleton: true, strictVersion: true, requiredVersion: 'auto'},
          "@angular/material/core": { singleton: true, strictVersion: true, requiredVersion: 'auto'},
          "@angular/material/form-field": { singleton: true, strictVersion: true, requiredVersion: 'auto'},
          "@angular/forms": { singleton: true, strictVersion: true, requiredVersion: 'auto'},
          "@angular/material/datepicker": { singleton: true, strictVersion: true, requiredVersion: 'auto'},
          "@angular/material/menu": { singleton: true, strictVersion: true, requiredVersion: 'auto'},
          "ngx-mat-timepicker": { singleton: true, strictVersion: true, requiredVersion: 'auto'},
          "@angular/platform-browser/animations": { singleton: true, strictVersion: true, requiredVersion: 'auto'},

          ...sharedMappings.getDescriptors()
        })

    }),
    sharedMappings.getPlugin()
  ],
};
