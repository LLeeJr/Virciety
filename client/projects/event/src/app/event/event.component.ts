import {Component, OnInit} from '@angular/core';
import {formatDate} from "@angular/common";
import {KeycloakService} from "keycloak-angular";
import {GQLService} from "../service/gql.service";
import {Event} from "../model/event";
import {AbstractControl, FormControl, FormGroup, ValidationErrors, ValidatorFn, Validators} from "@angular/forms";

@Component({
  selector: 'app-event',
  templateUrl: './event.component.html',
  styleUrls: ['./event.component.scss']
})
export class EventComponent implements OnInit {

  events: Event[] = [];
  currDate: string = formatDate(new Date(), 'fullDate', 'en-GB');

  description: string;
  location: string;
  username: string;
  checked: boolean;

  title: FormControl = new FormControl('', [Validators.required, emptyTextValidator('')]);
  startTime: FormControl = new FormControl('', Validators.required);
  endTime: FormControl = new FormControl('', Validators.required);
  range = new FormGroup({
    start: new FormControl('', Validators.required),
    end: new FormControl('', Validators.required),
  });

  constructor(private keycloak: KeycloakService,
              private gqlService: GQLService) { }

  async ngOnInit(): Promise<void> {
    await this.keycloak.isLoggedIn().then(loggedIn => {
      if (loggedIn) {
        this.keycloak.loadUserProfile().then(() => {
          this.username = this.keycloak.getUsername();
          this.gqlService.getEvents().subscribe(({data}: any) => {
            console.log(data);
            const unsorted: Event[] = [];

            for (let getEvent of data.getEvents) {
              const event: Event = new Event(getEvent);

              unsorted.push(event);
            }

            this.events = unsorted.sort((a, b) => +new Date(b.startDate) - +new Date(a.startDate))
          }, (error: any) => {
            console.error('there was an error sending the getEvents-query', error);
          });
        })
      } else {
        this.keycloak.login();
      }
    });

    /*this.gqlService.getEvents().subscribe(({data}: any) => {
      // console.log(data);
      const unsorted: Event[] = [];

      for (let getEvent of data.getEvents) {
        const event: Event = new Event(getEvent);

        unsorted.push(event);
      }

      this.events = unsorted.sort((a, b) => +new Date(b.startDate) - +new Date(a.startDate))
    }, (error: any) => {
      console.error('there was an error sending the getEvents-query', error);
    });*/
  }

  // TODO create dialog for this
  createEvent() {
    if (this.checkFields()) {
      return;
    }

    let startDate = formatDate(this.range.controls.start.value, 'fullDate', 'en-GB');
    let endDate = formatDate(this.range.controls.end.value, 'fullDate', 'en-GB');
    if (this.checked) {
      let shortDate = formatDate(startDate, 'short', 'en-GB');
      let split = shortDate.split(',', 1);
      startDate = split[0] + ', ' + this.startTime.value;

      shortDate = formatDate(endDate, 'short', 'en-GB');
      split = shortDate.split(',', 1);
      endDate = split[0] + ', ' + this.endTime.value;
    }


    this.gqlService
      .createEvent(this.title.value, this.username, startDate, endDate, this.location, this.description)
      .subscribe(({data}: any) => {
        console.log(data);
        this.events = [...this.events, data.createEvent]
      });
  }

  checkFields(): boolean {
    let required = false;

    if (!this.title.value) {
      this.title.markAsTouched();
      this.title.setErrors({required: true});
      required = true;
    }
    if (!this.range.controls.start.value) {
      this.range.controls.start.markAsTouched();
      this.range.controls.start.setErrors({required: true});
      required = true;
    }
    if (!this.range.controls.end.value) {
      this.range.controls.end.markAsTouched();
      this.range.controls.end.setErrors({required: true});
      required = true;
    }
    if (this.checked && !this.startTime.value) {
      this.startTime.markAsTouched();
      this.startTime.setErrors({required: true});
      required = true;
    }
    if (this.checked && !this.endTime.value) {
      this.endTime.markAsTouched();
      this.endTime.setErrors({required: true});
      required = true;
    }

    return required;
  }

  gotToMaps(location: string) {
    window.open('https://www.google.com/maps/search/?api=1&query=' + encodeURIComponent(location))
  }

  openEditDialog() {

  }

  checkTimes() {
    let startDate = formatDate(this.range.controls.start.value, 'fullDate', 'en-GB');
    let endDate = formatDate(this.range.controls.end.value, 'fullDate', 'en-GB');

    let startTime: string = this.startTime.value;
    let endTime: string = this.endTime.value;

    // compare dates and times
    if (startDate && endDate && startDate === endDate && startTime && endTime) {
      if (startTime.endsWith('PM') && endTime.endsWith('AM')) {
        this.endTime.setErrors({wrongTime: true});
      } else if (startTime.endsWith('AM') && endTime.endsWith('AM') || startTime.endsWith('PM') && endTime.endsWith('PM')) {
        let startHour: number = parseInt(startTime.split(':')[0]);
        let endHour: number = parseInt(endTime.split(':')[0]);

        if (startHour >= endHour) {
          this.endTime.setErrors({wrongTime: true});
        } else {
          this.endTime.setErrors(null);
        }
      } else {
        this.endTime.setErrors(null);
      }
    }
  }

  checkDate() {
    let startDate = formatDate(this.range.controls.start.value, 'fullDate', 'en-GB');

    if (+new Date(startDate) - +new Date(this.currDate) < 0) {
      this.range.controls.start.markAsTouched();
      this.range.controls.start.setErrors({wrongDate: true});
    } else {
      this.range.controls.start.setErrors(null);
    }
  }
}

export function emptyTextValidator(_: string): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    const emptyText = control.value !== '' && control.value.trim().length === 0;
    return emptyText ? {emptyText: true} : null;
  }
}
