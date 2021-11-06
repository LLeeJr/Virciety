import {Component, OnInit} from '@angular/core';
import {MfeOptions} from "./mfe/mfe";
import {LookupService} from "./mfe/lookup.service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'shell';
  mfes: MfeOptions[] = [];

  activeMfe: MfeOptions;

  constructor(
    private lookupService: LookupService
  ) { }

  async ngOnInit(): Promise<void> {
    this.mfes = await this.lookupService.lookup();
  }

  activate(mfe: MfeOptions): void {
    this.activeMfe = mfe;
  }
}
