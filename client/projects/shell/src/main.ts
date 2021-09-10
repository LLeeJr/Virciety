import {loadRemoteEntry} from "@angular-architects/module-federation";

Promise.all([
  loadRemoteEntry('http://localhost:5001/chatremoteEntry.js', 'chat'),
])
  .catch(err => console.error('Error loading remote entries', err))
  .then(() => import('./bootstrap'))
  .catch(err => console.error(err));
