<script lang=ts>
    export let docID: string;

    type document = {
        id: string,
        name: string,
        content: string,
        createdAt: string,
    }

    let data: document;

    $: {
        fetchDocument(docID);
    };

    async function fetchDocument(docID: string) {
        if (docID === undefined) {
            return;
        }
        let url = 'http://localhost:5500/document/' + docID;
        const res = await fetch(url);
        data = await res.json();
    }

</script>
{ #if data === undefined }
    <p>Loading...</p>
{:else}
    <h2>{data.name}</h2>
    <p >Created at: {data.createdAt}</p>
    <div>
        {data.content}
    </div>
{/if}
