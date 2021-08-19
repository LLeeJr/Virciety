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
    uniqueName: "commentMfe",
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
        name: "comment_mfe",
        filename: "comment_mfe_remoteEntry.js",
        exposes: {
            './CommentModule': './projects/comment-mfe/src/app/comment/comment.module.ts',
        },


        // For hosts (please adjust)
        // remotes: {
        //     "host": "host@http://localhost:4200/remoteEntry.js",
        //     "mfe1": "mfe1@http://localhost:5000/remoteEntry.js",
        //     "post": "post@http://localhost:5001/remoteEntry.js",
        //     "postMfe": "postMfe@http://localhost:5001/remoteEntry.js",

        // },

        shared: share({
          "@angular/core": { singleton: true, strictVersion: true, requiredVersion: 'auto' },
          "@angular/common": { singleton: true, strictVersion: true, requiredVersion: 'auto' },
          "@angular/common/http": { singleton: true, strictVersion: true, requiredVersion: 'auto' },
          "@angular/router": { singleton: true, strictVersion: true, requiredVersion: 'auto' },

          ...sharedMappings.getDescriptors()
        })

    }),
    sharedMappings.getPlugin()
  ],
};
