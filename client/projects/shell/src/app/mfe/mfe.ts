import { LoadRemoteModuleOptions } from '@angular-architects/module-federation';

export type MfeOptions = LoadRemoteModuleOptions & {
    displayName: string;
    componentName: string;
};
