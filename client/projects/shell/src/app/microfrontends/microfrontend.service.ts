import { Injectable } from "@angular/core";
import { Microfrontend } from "./microfrontend.model";

@Injectable({ providedIn: "root" })
export class MicrofrontendService {
  microfrontends: Microfrontend[];

  constructor() {}

  /*
   * Initialize is called on app startup to load the initial list of
   * remote microfrontends and configure them within the router
   */
  initialise(): Promise<void> {
    return new Promise<void>((resolve, _) => {
      this.microfrontends = this.loadConfig();
      resolve();
    });
  }

  /*
   * This is just an hardcoded list of remote microfrontends, but could easily be updated
   * to load the config from a database or external file
   */
  loadConfig(): Microfrontend[] {
    return [
      {
        // For Loading
        remoteEntry: "http://localhost:5002/remoteEntry.js",
        remoteName: "post",
        exposedModule: "CreatePostModule",

        // For Routing, enabling us to ngFor over the microfrontends and dynamically create links for the routes
        displayName: "CreatePostModule",
        routePath: "create_post",
        ngModuleName: "CreatePostModule",
      },
    ];
  }
}
