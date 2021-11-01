'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class NoShardingWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);


            // const accountID = (this.workerIndex==0) ? this.roundArguments.prefixA : this.roundArguments.prefixB;
            // console.log(`Worker ${this.workerIndex}: InitLedger ${accountID}`);
            // const request = {
            //     contractId: this.roundArguments.contractId,
            //     contractFunction: 'InitLedger',
            //     invokerIdentity: 'User1',
            //     contractArguments: [accountID, this.roundArguments.accounts],
            //     readOnly: false
            // };
            //
            // await this.sutAdapter.sendRequests(request);


    }

    async submitTransaction() {
        const randomId1 = Math.floor(Math.random()*this.roundArguments.accounts) + 1;
        const randomId2 = Math.floor(Math.random()*this.roundArguments.accounts) + 1;
        const myArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'SendAmount',
            invokerIdentity: 'User1',
            contractArguments: [`${this.roundArguments.prefixA}${randomId1}`, `${this.roundArguments.prefixB}${randomId2}`, 1],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(myArgs);
    }

    async cleanupWorkloadModule() {
        // for (let i=0; i<this.roundArguments.accounts; i++) {
        //     const accountID = `${this.roundArguments.prefixA}${i}`;
        //     console.log(`Worker ${this.workerIndex}: Deleting account ${accountID}`);
        //     const request = {
        //         contractId: this.roundArguments.contractId,
        //         contractFunction: 'DeleteAccount',
        //         invokerIdentity: 'User1',
        //         contractArguments: [accountID],
        //         readOnly: false
        //     };
        //
        //     await this.sutAdapter.sendRequests(request);
        // }
    }
}

function createWorkloadModule() {
    return new NoShardingWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
