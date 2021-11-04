import {loadRemoteEntry} from "@angular-architects/module-federation";

Promise.all([
  loadRemoteEntry('http://localhost:5001/chatremoteEntry.js', 'chat'),
  loadRemoteEntry('http://localhost:5002/postremoteEntry.js', 'post'),
  loadRemoteEntry('http://localhost:5003/userremoteEntry.js', 'user'),
])
  .catch(err => console.error('Error loading remote entries', err))
  .then(() => import('./bootstrap'))
  .catch(err => console.error(err));
