const ModuleFederationPlugin = require("webpack/lib/container/ModuleFederationPlugin");
const mf = require("@angular-architects/module-federation/webpack");
const path = require("path");
const share = mf.share;

const sharedMappings = new mf.SharedMappings();
sharedMappings.register(
  path.join(__dirname, '../../tsconfig.json'),
  ['data-lib']);

module.exports = {
  output: {
    uniqueName: "postMfe",
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
        name: "post_mfe",
        filename: "post_mfe_remoteEntry.js",
        exposes: {
            './PostModule': './projects/post-mfe/src/app/post/post.module.ts',
        },

        // For hosts (please adjust)
        // remotes: {
        //     "host": "host@http://localhost:4200/remoteEntry.js",
        //     "mfe1": "mfe1@http://localhost:5000/remoteEntry.js",
        //     "post": "post@http://localhost:5001/remoteEntry.js",

        // },

        shared: share({
          "@angular/core": { singleton: true, strictVersion: true, requiredVersion: 'auto' },
          "@angular/common": { singleton: true, strictVersion: true, requiredVersion: 'auto' },
          "@angular/common/http": { singleton: true, strictVersion: true, requiredVersion: 'auto' },
          "@angular/router": { singleton: true, strictVersion: true, requiredVersion: 'auto' },
          "@apollo/client": { singleton: true, strictVersion: true, requiredVersion: 'auto' },
          "apollo-angular": { singleton: true, strictVersion: true, requiredVersion: 'auto' },

          ...sharedMappings.getDescriptors()
        })

    }),
    sharedMappings.getPlugin()
  ],
};
