import {Values} from "./values";

export namespace Workflows {
    import CandidateDetails = Values.CandidateDetails;

    export interface BackgroundCheckWorkflowInput {
        email: string,
        tier: string,
    }
    export interface BackgroundCheckWorkflowResult {
        email: string,
        tier: string,
        accepted: boolean,
        candidateDetails: CandidateDetails,
        ssnTrace?: SSNTraceWorkflowResult,
        searchResults: Map<string,any>,
        searchErrors: Map<string, any>,
    }

    export interface SSNTraceWorkflowResult {

    }
}

export namespace Commands {

}
