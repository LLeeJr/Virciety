import {Component, Inject, OnInit} from '@angular/core';
import {formatDate} from "@angular/common";
import {AbstractControl, FormControl, FormGroup, ValidationErrors, ValidatorFn, Validators} from "@angular/forms";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import {CreateEventData, CreateEventDialogData} from "../model/dialog.data";

@Component({
  selector: 'app-create-event',
  templateUrl: './create-event.component.html',
  styleUrls: ['./create-event.component.scss']
})
export class CreateEventComponent implements OnInit {

  currDate: string = formatDate(new Date(), 'fullDate', 'en-GB');
  delete: CreateEventData = {
    event: this.data.event,
    remove: true,
  }

  description: string;
  location: string;
  checked: boolean;

  title: FormControl = new FormControl('', [Validators.required, emptyTextValidator('')]);
  startTime: FormControl = new FormControl('', Validators.required);
  endTime: FormControl = new FormControl('', Validators.required);
  range = new FormGroup({
    startDate: new FormControl('', Validators.required),
    endDate: new FormControl('', Validators.required),
  });

  constructor(@Inject(MAT_DIALOG_DATA) public data: CreateEventDialogData,
              private dialogRef: MatDialogRef<CreateEventComponent>) {
    // fill the dialog with event data if in edit mode
    if (data.editMode && data.event) {
      this.title.setValue(data.event.title);
      this.range.controls.startDate.setValue(new Date(data.event.startDate));
      this.range.controls.endDate.setValue(new Date(data.event.endDate));
      this.location = data.event.location;
      this.description = data.event.description;
      this.checked = !!(data.event.startTime && data.event.endTime);

      if (this.checked) {
        this.startTime.setValue(data.event.startTime);
        this.endTime.setValue(data.event.endTime);
      }
    }
  }

  ngOnInit(): void {
  }

  // check if all required fields have a value, if so close dialog and do an action. If not, mark the required fields
  submit(): void {
    // if in edit mode, update event data
    if (this.data.editMode) {
      if (this.data.event) {
        this.data.event.title = this.title.value;
        this.data.event.location = this.location;
        this.data.event.description = this.description;
        this.data.event.startDate = this.range.controls.startDate.value;
        this.data.event.endDate = this.range.controls.endDate.value;

        if (this.checked) {
          this.data.event.startTime = this.startTime.value;
          this.data.event.endTime = this.endTime.value;
        } else {
          this.data.event.startTime = null;
          this.data.event.endTime = null;
        }
      }

      let output: CreateEventData = {
        event: this.data.event,
        remove: false
      }

      this.dialogRef.close(output);
    // if not in editMode create output data, so new event can be created
    } else {
      let output: CreateEventData = {
        event: {
          description: this.description,
          location: this.location,
          startDate: this.range.controls.startDate.value,
          endDate: this.range.controls.endDate.value,
          startTime: this.startTime.value === '' ? null : this.startTime.value,
          endTime: this.endTime.value === '' ? null : this.endTime.value,
          title: this.title.value,
        },
        remove: false
      }
      this.dialogRef.close(output);
    }
  }

  // check if input times are valid
  checkTimes(): boolean {
    if (this.range.controls.startDate.value && this.range.controls.endDate.value) {
      let startDate = formatDate(this.range.controls.startDate.value, 'fullDate', 'en-GB');
      let endDate = formatDate(this.range.controls.endDate.value, 'fullDate', 'en-GB');

      let startTime: string = this.startTime.value;
      let endTime: string = this.endTime.value;

      // check if datetime is valid
      if (startDate && endDate && startDate === endDate && startTime && endTime) {
        if (startTime.endsWith('PM') && endTime.endsWith('AM')) {
          this.endTime.setErrors({wrongTime: true});
          return false;
        } else if (startTime.endsWith('AM') && endTime.endsWith('AM') || startTime.endsWith('PM') && endTime.endsWith('PM')) {
          let startHour: number = parseInt(startTime.split(':')[0]);
          let endHour: number = parseInt(endTime.split(':')[0]);

          if ((startHour !== 12 && startHour >= endHour) || endHour === 12) {
            this.endTime.setErrors({wrongTime: true});
            return false;
          } else {
            this.endTime.setErrors(null);
            return true;
          }
        } else {
          this.endTime.setErrors(null);
          return true;
        }
      } else {
        if (this.checked && startTime && endTime) {
          return true;
        } else return !(this.checked && !(startTime && endTime));
      }
    } else {
      return false;
    }
  }

  // check that input startDate isn't in the past
  checkDate(): boolean {
    let startDate = formatDate(this.range.controls.startDate.value, 'fullDate', 'en-GB');

    if (+new Date(startDate) - +new Date(this.currDate) < 0) {
      this.range.controls.startDate.markAsTouched();
      this.range.controls.startDate.setErrors({pastDate: true});
      return false;
    } else {
      this.range.controls.startDate.setErrors(null);
      return true;
    }
  }

  fieldsValid() {
    return this.title.valid && this.checkDate() && this.checkTimes();
  }
}

export function emptyTextValidator(_: string): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    const emptyText = control.value !== '' && control.value.trim().length === 0;
    return emptyText ? {emptyText: true} : null;
  }
}

