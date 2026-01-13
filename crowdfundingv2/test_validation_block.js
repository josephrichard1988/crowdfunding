
// const AUTH_API = 'http://localhost:3001/api/auth';
// const NETWORK_API = 'http://localhost:4000/api/startup';

// async function test() {
//     try {
//         console.log('Using Node build-in fetch...');

//         // 1. Create Startup
//         console.log('Creating startup...');
//         const startupRes = await fetch(`${NETWORK_API}/startups`, {
//             method: 'POST',
//             headers: { 'Content-Type': 'application/json' },
//             body: JSON.stringify({ name: 'NoValStartupNode3', description: 'Desc', ownerId: 'NOVAL_OWNER_NODE3' })
//         });
//         const startupData = await startupRes.json();
//         console.log('Startup Response:', JSON.stringify(startupData));
//         const startupId = startupData.data?.startupId;

//         if (!startupId) {
//             throw new Error('Failed to create startup: ' + JSON.stringify(startupData));
//         }
//         console.log('Startup ID:', startupId);

//         console.log('Waiting 15s for propagation...');
//         await new Promise(r => setTimeout(r, 15000));

//         // 2. Create Campaign
//         console.log('Creating campaign...');
//         const campRes = await fetch(`${NETWORK_API}/campaigns`, {
//             method: 'POST',
//             headers: { 'Content-Type': 'application/json' },
//             body: JSON.stringify({
//                 startupId,
//                 projectName: 'NoValTestNode3',
//                 description: 'Desc',
//                 goalAmount: 1000,
//                 category: 'Technology',
//                 ownerId: 'NOVAL_OWNER_NODE3', // This parameter might be redundant or wrong, checking logs
//                 deadline: '2026-06-01',
//                 currency: 'USD',
//                 fundingMonth: 1,
//                 fundingDay: 1,
//                 fundingYear: 2026
//             })
//         });
//         const campData = await campRes.json();
//         console.log('Campaign Response:', JSON.stringify(campData));
//         const campaignId = campData.data?.campaignId;

//         if (!campaignId) {
//             throw new Error('Failed to create campaign: ' + JSON.stringify(campData));
//         }

//         // 3. Login
//         console.log('Logging in...');
//         const loginRes = await fetch(`${AUTH_API}/login`, {
//             method: 'POST',
//             headers: { 'Content-Type': 'application/json' },
//             body: JSON.stringify({ email: 'startup@example.com', password: 'password123', role: 'STARTUP' })
//         });
//         const loginData = await loginRes.json();
//         const token = loginData.token;
//         console.log('Token acquired');

//         // 4. Submit for Validation (Expect Failure)
//         console.log('Submitting for validation...');
//         const submitRes = await fetch(`${NETWORK_API}/campaigns/${campaignId}/submit-validation`, {
//             method: 'POST',
//             headers: {
//                 'Content-Type': 'application/json',
//                 'Authorization': `Bearer ${token}`
//             },
//             body: JSON.stringify({
//                 authToken: token,
//                 documents: [],
//                 notes: 'Should fail'
//             })
//         });

//         const submitData = await submitRes.json();

//         if (submitRes.status === 400 || submitRes.status === 503) {
//             console.log('SUCCESS: Submission failed as expected with status', submitRes.status);
//             console.log('Error message:', submitData.error);
//         } else {
//             console.log('FAILURE: Submission succeeded with status', submitRes.status);
//             console.log('Response:', submitData);
//         }

//     } catch (e) {
//         console.error('Test Failed:', e.message);
//     }
// }

// test();
