import {Event} from "./event";

export interface CreateEventDialogData {
  editMode: boolean;
  event: Event | null;
}

export interface CreateEventData {
  event: Event | {
    description: string,
    location: string,
    startDate: string,
    endDate: string,
    startTime: string | null,
    endTime: string | null,
    title: string,
  } | null;
  remove: boolean;
}
