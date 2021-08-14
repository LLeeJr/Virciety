import {loadRemoteEntry} from "@angular-architects/module-federation";

Promise.all([
  loadRemoteEntry('http://localhost:5000/mfe1remoteEntry.js', 'mfe1'),
  loadRemoteEntry('http://localhost:5001/post_mfe_remoteEntry.js', 'post_mfe'),
  loadRemoteEntry('http://localhost:5002/comment_mfe_remoteEntry.js', 'comment_mfe')
])
  .catch(err => console.error('Error loading remote entries', err))
  .then(() => import('./bootstrap'))
  .catch(err => console.error(err));
