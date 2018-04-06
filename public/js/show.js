$(document).ready(function() {
    seasonTable = $('#seasons').DataTable({
        "ordering": false,
        "info":     false,
        "lengthChange": false,
        "pageLength": 5,
        "sDom": '<"top"i>rt<"bottom"lp><"clear">',
    });
    timeslotTable = $('#timeslots').DataTable({
        "ordering": false,
        "info":     false,
        "lengthChange": false,
        "sDom": '<"top"i>rt<"bottom"lp><"clear">',
        "pageLength": 5,
        "columnDefs": [
        {
            "targets": [ 0 ],
            "visible": false,
            "searchable": true
        },
        {
            "targets": [ 1,2 ],
            "searchable": false
        },

        ]
    });
    $('.timeslotFilterSeason').on( 'click', function () {
        $('.timeslotFilterSeason').removeClass("active");
        $(this).addClass("active");
        $('#timeslotLatest').hide();
        $('#timeslotSeasonTitle').text("Season " + $(this).find('td:first-of-type').html())
        timeslotTable
            .columns( 0 )
            .search( $(this).data("seasonid") )
            .draw();

        $('#timeslotList').show();

        $('#timeslotList tr').on("click", function () {
            window.location.href = $(this).find("a").attr("href");
        });
    } );
} );
