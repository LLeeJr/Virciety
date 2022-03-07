const ModuleFederationPlugin = require("webpack/lib/container/ModuleFederationPlugin");
const mf = require("@angular-architects/module-federation/webpack");
const path = require("path");
const share = mf.share;

const sharedMappings = new mf.SharedMappings();
sharedMappings.register(
  path.join(__dirname, '../../tsconfig.json'),
  ['auth-lib']);

module.exports = {
  output: {
    uniqueName: "notification",
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
        name: "notification",
        filename: "remoteEntry.js",
        exposes: {
          './NotificationModule': './projects/notification/src/app/notification/notification.module.ts',
          NotificationModule: './projects/notification/src/app/notification/notification.module.ts',
        },

        // For hosts (please adjust)
        // remotes: {
        //     "shell": "shell@http://localhost:4200/remoteEntry.js",
        //     "chat": "chat@http://localhost:5001/remoteEntry.js",
        //     "post": "post@http://localhost:5002/remoteEntry.js",
        //     "user": "user@http://localhost:5003/remoteEntry.js",
        //     "profileViewer": "profileViewer@http://localhost:5005/remoteEntry.js",

        // },

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

          ...sharedMappings.getDescriptors()
        })

    }),
    sharedMappings.getPlugin()
  ],
};
